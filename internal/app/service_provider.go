package app

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/client"
	"github.com/balobas/auth_service/internal/client/pg"
	"github.com/balobas/auth_service/internal/config"
	deliveryGrpc "github.com/balobas/auth_service/internal/delivery/grpc"
	"github.com/balobas/auth_service/internal/manager/transaction"
	repositoryPostgres "github.com/balobas/auth_service/internal/repository/postgres"
	"github.com/balobas/auth_service/internal/shutdown"
)

type serviceProvider struct {
	pgConfig   *config.ConfigPG
	grpcConfig *config.ConfigGRPC

	pgClient  client.ClientDB
	usersRepo *repositoryPostgres.Repository

	txManager *transaction.Manager

	usersService   *usersService.UsersService
	authServerGrpc *deliveryGrpc.AuthServerGrpc
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) PgConfig() *config.ConfigPG {
	if sp.pgConfig == nil {
		sp.pgConfig = config.NewConfigPG()
	}
	return sp.pgConfig
}

func (sp *serviceProvider) GrpcConfig() *config.ConfigGRPC {
	if sp.grpcConfig == nil {
		cfg, err := config.NewConfigGRPC()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		sp.grpcConfig = cfg
	}
	return sp.grpcConfig
}

func (sp *serviceProvider) PgClient(ctx context.Context) client.ClientDB {
	if sp.pgClient == nil {
		client, err := pg.NewClient(ctx, sp.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create pgClient: %v", err)
		}

		shutdown.Add(client.Close)

		sp.pgClient = client
	}
	return sp.pgClient
}

func (sp *serviceProvider) UsersRepo(ctx context.Context) *repositoryPostgres.Repository {
	if sp.usersRepo == nil {
		sp.usersRepo = repositoryPostgres.New(sp.PgClient(ctx))
	}
	return sp.usersRepo
}

func (sp *serviceProvider) TxManager() *transaction.Manager {
	if sp.txManager == nil {
		sp.txManager = transaction.NewTxManager()
	}
	return sp.txManager
}

func (sp *serviceProvider) UsersService(ctx context.Context) *usersService.UsersService {
	if sp.usersService == nil {
		sp.usersService = usersService.New(sp.UsersRepo(ctx))
	}
	return sp.usersService
}

func (sp *serviceProvider) AuthServerGrpc(ctx context.Context) *deliveryGrpc.AuthServerGrpc {
	if sp.authServerGrpc == nil {
		sp.authServerGrpc = deliveryGrpc.NewAuthServerGRPC(
			sp.GrpcConfig(),
			sp.UsersService(ctx),
		)
	}
	return sp.authServerGrpc
}
