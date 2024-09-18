package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/balobas/auth_service/internal/config"
	"github.com/balobas/auth_service/internal/shutdown"
	"github.com/balobas/auth_service/pkg/auth_v1"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		shutdown.CloseAll(shutdownCtx)
	}()

	if err := a.initDeps(ctx); err != nil {
		return errors.Wrap(err, "failed to init app deps")
	}

	a.runGrpcServer(ctx)
	a.runVerificationWorker(ctx)

	<-ctx.Done()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initEnv,
		a.initServiceProvider,
		a.initGrpcServer,
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

func (a *App) initGrpcServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(a.serviceProvider.AuthServerGrpc(ctx).UnaryAuthInterceptor()),
	)
	reflection.Register(a.grpcServer)
	auth_v1.RegisterAuthServer(a.grpcServer, a.serviceProvider.AuthServerGrpc(ctx))

	return nil
}

func (a *App) runGrpcServer(ctx context.Context) {
	log.Printf("grpc server is running on %v\n", a.serviceProvider.GrpcConfig())

	lis, err := net.Listen("tcp", a.serviceProvider.GrpcConfig().Address())
	if err != nil {
		log.Fatalf("failed to listen tcp: %v", err)
	}

	go func() {
		done := make(chan struct{}, 1)
		go func() {
			err := a.grpcServer.Serve(lis)
			if err != nil {
				log.Default().Printf("grpc server cancelled with error: %v\n", err)
			} else {
				log.Default().Println("grpc server cancelled without errors")
			}
			done <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			log.Printf("grpc server cancelled, ctx.Done, error: %v", ctx.Err())
			return
		case <-done:
			log.Printf("grpc server cancelled")
		}
	}()
}

func (a *App) runVerificationWorker(ctx context.Context) {
	go a.serviceProvider.WorkerVerification(ctx).Run(ctx)
}
