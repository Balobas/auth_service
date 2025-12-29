package app

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/balobas/auth_service/internal/config"
	"github.com/balobas/auth_service/internal/shutdown"
	"github.com/balobas/auth_service/migrations"
	"github.com/balobas/auth_service/pkg/auth_v1"
	"github.com/balobas/sport_city_common/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server

	configPath string
}

func NewApp(configPath string) *App {
	return &App{configPath: configPath}
}

func (a *App) Run(ctx context.Context) error {
	done := make(chan struct{}, 2)

	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered panic in app: %v", err)
		}

		log.Printf("start shutdown")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		shutdown.CloseAll(shutdownCtx)
		close(done)
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := a.initDeps(ctx); err != nil {
		return errors.Wrap(err, "failed to init app deps")
	}

	a.migrateDb(ctx)
	a.runGrpcServer(ctx, done)
	a.runWorkers(ctx)
	a.runVerificationWorker(ctx)

	select {
	case <-ctx.Done():
	case <-done:
		log.Printf("one of the critical components stopped")
	}
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initEnv,
		a.initServiceProvider,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initEnv(_ context.Context) error {
	if err := config.Load(a.configPath); err != nil {
		return errors.Wrap(err, "failed to init env")
	}
	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) runGrpcServer(ctx context.Context, done chan<- struct{}) {
	log.Printf("grpc server is running on %v\n", a.serviceProvider.GrpcConfig())
	authGrpcServer := a.serviceProvider.AuthServerGrpc(ctx)

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			authGrpcServer.UnaryAuthInterceptor(),
			authGrpcServer.UnaryErrorsPostInterceptor(),
		),
	)
	reflection.Register(a.grpcServer)
	auth_v1.RegisterAuthServer(a.grpcServer, authGrpcServer)

	go func() {
		lis, err := net.Listen("tcp", a.serviceProvider.GrpcConfig().Address())
		if err != nil {
			log.Printf("failed to listen tcp: %v", err)
			done <- struct{}{}
			return
		}

		shutdown.Add(func(ctx context.Context) error {
			a.grpcServer.Stop()
			return nil
		})

		serverStopped := make(chan struct{})
		go func() {
			err := a.grpcServer.Serve(lis)
			if err != nil {
				log.Default().Printf("grpc server cancelled with error: %v\n", err)
			} else {
				log.Default().Println("grpc server cancelled without errors")
			}
			close(serverStopped)
		}()

		select {
		case <-ctx.Done():
			log.Printf("grpc server cancelled, ctx.Done, error: %v", ctx.Err())
			return
		case <-serverStopped:
			log.Printf("grpc server cancelled")
			done <- struct{}{}
		}
	}()
}

func (a *App) runVerificationWorker(ctx context.Context) {
	go a.serviceProvider.WorkerVerification(ctx).Run(ctx)
}

func (a *App) runWorkers(ctx context.Context) {
	sp := a.serviceProvider
	sp.RiverClient(ctx).BuildWorkers(
		ctx,
		sp.WorkerMqPublisher(ctx),
	)

	if err := a.serviceProvider.RiverClient(ctx).Start(ctx); err != nil {
		log := logger.From(ctx)
		log.Error().Msgf("failed to start river workers: %v", err)
		panic(err)
	}
}

func (a *App) migrateDb(ctx context.Context) {
	if err := a.serviceProvider.PgClient(ctx).Migrate(ctx, migrations.Files); err != nil {
		panic(errors.Wrap(err, "failed to migrate db"))
	}
}
