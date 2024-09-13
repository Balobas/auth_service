package app

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/client"
	"github.com/balobas/auth_service/internal/client/email"
	"github.com/balobas/auth_service/internal/client/pg"
	"github.com/balobas/auth_service/internal/config"
	deliveryGrpc "github.com/balobas/auth_service/internal/delivery/grpc"
	jwtManager "github.com/balobas/auth_service/internal/manager/jwt"
	"github.com/balobas/auth_service/internal/manager/transaction"
	repositoryKeys "github.com/balobas/auth_service/internal/repository/keys"
	repositoryConfig "github.com/balobas/auth_service/internal/repository/postgres/config"
	repositoryCredentials "github.com/balobas/auth_service/internal/repository/postgres/credentials"
	repositoryPermissions "github.com/balobas/auth_service/internal/repository/postgres/permissions"
	sessionRepository "github.com/balobas/auth_service/internal/repository/postgres/session"
	repositoryUsers "github.com/balobas/auth_service/internal/repository/postgres/users"
	repositoryVerification "github.com/balobas/auth_service/internal/repository/postgres/verification"
	"github.com/balobas/auth_service/internal/shutdown"
	useCaseAuth "github.com/balobas/auth_service/internal/usecase/auth"
	useCaseConfig "github.com/balobas/auth_service/internal/usecase/config"
	useCaseCredentials "github.com/balobas/auth_service/internal/usecase/credentials"
	useCaseUsers "github.com/balobas/auth_service/internal/usecase/users"
	useCaseVerification "github.com/balobas/auth_service/internal/usecase/verification"
	workerVerification "github.com/balobas/auth_service/internal/worker/verification"
)

type serviceProvider struct {
	pgConfig      *config.ConfigPG
	grpcConfig    *config.ConfigGRPC
	serviceConfig *config.ServiceConfig

	pgClient    client.ClientDB
	emailClient *email.SmtpClient

	keysRepository         *repositoryKeys.KeysRepository
	usersRepository        *repositoryUsers.UsersRepository
	permissionsRepository  *repositoryPermissions.PermissionsRepository
	credentialsRepository  *repositoryCredentials.CredentialsRepository
	sessionsRepository     *sessionRepository.SessionRepository
	verificationRepository *repositoryVerification.VerificationRepository
	configRepository       *repositoryConfig.ConfigRepository

	txManager  *transaction.Manager
	jwtManager *jwtManager.JwtManager

	useCaseConfig       *useCaseConfig.UseCaseConfig
	useCaseUsers        *useCaseUsers.UseCaseUsers
	useCaseCredentials  *useCaseCredentials.UseCaseCredentials
	useCaseVerification *useCaseVerification.UseCaseVerification
	useCaseAuth         *useCaseAuth.UseCaseAuth

	workerVerification *workerVerification.Worker

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
