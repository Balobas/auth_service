package app

import (
	"context"
	"log"
	"net"
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
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		shutdown.CloseAll(shutdownCtx)
	}()

	if err := a.initDeps(ctx); err != nil {
		return errors.Wrap(err, "failed to init app deps")
	}

	return a.runGrpcServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
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

func (a *App) initConfig(_ context.Context) error {
	if err := config.Load(a.configPath); err != nil {
		return errors.Wrap(err, "failed to init config")
	}

	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	reflection.Register(a.grpcServer)
	auth_v1.RegisterAuthServer(a.grpcServer, a.serviceProvider.AuthServerGrpc(ctx))
	return nil
}

func (a *App) runGrpcServer() error {
	log.Printf("grpc server is running on %v\n", a.serviceProvider.GrpcConfig())

	lis, err := net.Listen("tcp", a.serviceProvider.GrpcConfig().Address())
	if err != nil {
		return errors.Wrap(err, "failed to listen tcp")
	}

	return a.grpcServer.Serve(lis)
}
