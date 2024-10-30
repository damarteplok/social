package main

import (
	"time"

	"github.com/damarteplok/social/internal/auth"
	"github.com/damarteplok/social/internal/db"
	"github.com/damarteplok/social/internal/env"
	"github.com/damarteplok/social/internal/mailer"
	"github.com/damarteplok/social/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gohpers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:        env.Envs.Addr,
		frontendURL: env.Envs.FrontendURL,
		db: dbConfig{
			addr:         env.Envs.DbAddr,
			maxOpenConns: env.Envs.MaxOpenConns,
			maxIdleConns: env.Envs.MaxIdleConns,
			maxIdleTime:  env.Envs.MaxIdleTime,
		},
		env:    env.Envs.ENV,
		apiURL: env.Envs.ApiUrl,
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.Envs.MailerFromEmail,
			sendgrid: sendGridConfig{
				apiKey: env.Envs.MailerApiKey,
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.Envs.AdminUser,
				pass: env.Envs.AdminPass,
			},
			token: tokenConfig{
				secret: env.Envs.JwtSecret,
				exp:    time.Hour * 24 * 3,
				iss:    env.Envs.JwtIss,
			},
		},
		camunda: camundaConfig{
			zeebeAddr:          env.Envs.ZeebeAddr,
			zeebeClientId:      env.Envs.ZeebeClientID,
			zeebeClientSecret:  env.Envs.ZeebeClientSecret,
			zeebeAuthServerUrl: env.Envs.ZeebeAuthServerUrl,
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(db)

	mailer := mailer.NewSendgrid(cfg.mail.sendgrid.apiKey, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
