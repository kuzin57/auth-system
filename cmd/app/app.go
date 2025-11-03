package app

import (
	"os"

	"github.com/kuzin57/auth-system/internal/broker"
	"github.com/kuzin57/auth-system/internal/config"
	"github.com/kuzin57/auth-system/internal/handlers"
	"github.com/kuzin57/auth-system/internal/repositories"
	usersrepo "github.com/kuzin57/auth-system/internal/repositories/users"
	"github.com/kuzin57/auth-system/internal/services/email"
	"github.com/kuzin57/auth-system/internal/services/token"
	"github.com/kuzin57/auth-system/internal/services/users"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func Create(confPath, secretsPath string) fx.Option {
	var (
		conf    = mustLoadConfig(confPath)
		secrets = mustLoadSecrets(secretsPath)
	)

	return fx.Options(
		fx.Supply(conf, secrets),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			zap.NewProduction,
			repositories.NewPgDriver,
			usersrepo.NewRepository,
			users.NewService,
			handlers.NewSendRegistrationLinkHandler,
			handlers.NewRegistrationPageHandler,
			handlers.NewListUsersHandler,
			handlers.NewRegisterHandler,
			handlers.NewServer,
			broker.NewMessageBroker,
			email.NewService,
			token.NewService,
		),
		fx.Invoke(func(server *handlers.Server) {}),
		fx.Invoke(func(broker *broker.MessageBroker) {}),
		fx.Invoke(func(emailService *email.Service) {}),
	)
}

func mustLoadConfig(path string) *config.Config {
	confContent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	conf := &config.Config{}

	if err := yaml.Unmarshal(confContent, conf); err != nil {
		panic(err)
	}

	return conf
}

func mustLoadSecrets(path string) *config.Secrets {
	secretsContent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	secrets := &config.Secrets{}

	if err := yaml.Unmarshal(secretsContent, secrets); err != nil {
		panic(err)
	}

	return secrets
}
