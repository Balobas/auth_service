package app

import (
	"context"
	"log"

	"github.com/balobas/auth_service/internal/client"
	"github.com/balobas/auth_service/internal/client/pg"
	"github.com/balobas/auth_service/internal/config"
	deliveryGrpc "github.com/balobas/auth_service/internal/delivery/grpc"
	jwtManager "github.com/balobas/auth_service/internal/manager/jwt"
	"github.com/balobas/auth_service/internal/manager/transaction"
	emailMock "github.com/balobas/auth_service/internal/mocks/email"
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
	emailClient EmailClient

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

func (sp *serviceProvider) ServiceConfig() *config.ServiceConfig {
	if sp.serviceConfig == nil {
		sp.serviceConfig = config.NewServiceConfig()
	}
	return sp.serviceConfig
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

func (sp *serviceProvider) EmailClient(ctx context.Context) EmailClient {
	if sp.emailClient == nil {
		// sp.emailClient = email.NewClient(ctx, sp.ServiceConfig())
		sp.emailClient = emailMock.NewClient(ctx, sp.ServiceConfig())
	}
	return sp.emailClient
}

func (sp *serviceProvider) KeysRepository(ctx context.Context) *repositoryKeys.KeysRepository {
	if sp.keysRepository == nil {
		sp.keysRepository = repositoryKeys.New()
	}
	return sp.keysRepository
}

func (sp *serviceProvider) UsersRepository(ctx context.Context) *repositoryUsers.UsersRepository {
	if sp.usersRepository == nil {
		sp.usersRepository = repositoryUsers.New(sp.PgClient(ctx))
	}
	return sp.usersRepository
}

func (sp *serviceProvider) PermissionsRepository(ctx context.Context) *repositoryPermissions.PermissionsRepository {
	if sp.permissionsRepository == nil {
		sp.permissionsRepository = repositoryPermissions.New(sp.PgClient(ctx))
	}
	return sp.permissionsRepository
}

func (sp *serviceProvider) CredentialsRepository(ctx context.Context) *repositoryCredentials.CredentialsRepository {
	if sp.credentialsRepository == nil {
		sp.credentialsRepository = repositoryCredentials.New(sp.PgClient(ctx))
	}
	return sp.credentialsRepository
}

func (sp *serviceProvider) SessionsRepository(ctx context.Context) *sessionRepository.SessionRepository {
	if sp.sessionsRepository == nil {
		sp.sessionsRepository = sessionRepository.New(sp.PgClient(ctx))
	}
	return sp.sessionsRepository
}

func (sp *serviceProvider) VerificationRepository(ctx context.Context) *repositoryVerification.VerificationRepository {
	if sp.verificationRepository == nil {
		sp.verificationRepository = repositoryVerification.New(sp.PgClient(ctx))
	}
	return sp.verificationRepository
}

func (sp *serviceProvider) ConfigRepository(ctx context.Context) *repositoryConfig.ConfigRepository {
	if sp.configRepository == nil {
		sp.configRepository = repositoryConfig.New(sp.PgClient(ctx))
	}
	return sp.configRepository
}

func (sp *serviceProvider) TxManager(ctx context.Context) *transaction.Manager {
	if sp.txManager == nil {
		sp.txManager = transaction.NewTxManager(sp.PgClient(ctx))
	}
	return sp.txManager
}

func (sp *serviceProvider) JwtManager(ctx context.Context) *jwtManager.JwtManager {
	if sp.jwtManager == nil {
		sp.jwtManager = jwtManager.New(sp.KeysRepository(ctx))
	}
	return sp.jwtManager
}

func (sp *serviceProvider) UseCaseConfig(ctx context.Context) *useCaseConfig.UseCaseConfig {
	if sp.useCaseConfig == nil {
		sp.useCaseConfig = useCaseConfig.New(sp.ServiceConfig(), sp.ConfigRepository(ctx))
	}
	return sp.useCaseConfig
}

func (sp *serviceProvider) initConfig(ctx context.Context) {

	if err := sp.UseCaseConfig(ctx).InitFromDB(ctx); err != nil {
		log.Fatal("failed to init service config")
	}
}

func (sp *serviceProvider) UseCaseUsers(ctx context.Context) *useCaseUsers.UseCaseUsers {
	if sp.useCaseUsers == nil {
		sp.useCaseUsers = useCaseUsers.New(
			sp.UsersRepository(ctx),
			sp.PermissionsRepository(ctx),
			sp.UseCaseVerification(ctx),
			sp.TxManager(ctx),
			sp.UseCaseCredentials(ctx),
		)
	}
	return sp.useCaseUsers
}

func (sp *serviceProvider) UseCaseCredentials(ctx context.Context) *useCaseCredentials.UseCaseCredentials {
	if sp.useCaseCredentials == nil {
		sp.useCaseCredentials = useCaseCredentials.New(
			sp.ServiceConfig(),
			sp.CredentialsRepository(ctx),
		)
	}
	return sp.useCaseCredentials
}

func (sp *serviceProvider) UseCaseVerification(ctx context.Context) *useCaseVerification.UseCaseVerification {
	if sp.useCaseVerification == nil {
		sp.useCaseVerification = useCaseVerification.New(
			sp.ServiceConfig(),
			sp.VerificationRepository(ctx),
			sp.PermissionsRepository(ctx),
			sp.TxManager(ctx),
		)
	}
	return sp.useCaseVerification
}

func (sp *serviceProvider) UseCaseAuth(ctx context.Context) *useCaseAuth.UseCaseAuth {
	if sp.useCaseAuth == nil {
		sp.useCaseAuth = useCaseAuth.New(
			sp.ServiceConfig(),
			sp.SessionsRepository(ctx),
			sp.PermissionsRepository(ctx),
			sp.UseCaseUsers(ctx),
			sp.UseCaseCredentials(ctx),
			sp.JwtManager(ctx),
			sp.TxManager(ctx),
		)
	}
	return sp.useCaseAuth
}

func (sp *serviceProvider) WorkerVerification(ctx context.Context) *workerVerification.Worker {
	if sp.workerVerification == nil {
		sp.workerVerification = workerVerification.New(
			sp.ServiceConfig(),
			sp.VerificationRepository(ctx),
			sp.EmailClient(ctx),
		)
	}
	return sp.workerVerification
}

func (sp *serviceProvider) AuthServerGrpc(ctx context.Context) *deliveryGrpc.AuthServerGrpc {
	if sp.authServerGrpc == nil {
		sp.initConfig(ctx)

		sp.authServerGrpc = deliveryGrpc.NewAuthServerGRPC(
			sp.ServiceConfig(),
			sp.UseCaseUsers(ctx),
			sp.UseCaseAuth(ctx),
			sp.UseCaseVerification(ctx),
		)
	}
	return sp.authServerGrpc
}

type EmailClient interface {
	SendEmail(receiverEmail string, body []byte) error
}
