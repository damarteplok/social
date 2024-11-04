package main

import (
	"expvar"
	"runtime"

	"github.com/damarteplok/social/internal/auth"
	"github.com/damarteplok/social/internal/db"
	"github.com/damarteplok/social/internal/env"
	"github.com/damarteplok/social/internal/mailer"
	"github.com/damarteplok/social/internal/ratelimiter"
	"github.com/damarteplok/social/internal/store"
	"github.com/damarteplok/social/internal/store/cache"
	"github.com/damarteplok/social/internal/zeebe"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			damarmunda API
//	@description	API for damarmunda, a camunda golang
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
		redisCfg: redisConfig{
			addr:    env.Envs.RedisAddr,
			pw:      env.Envs.RedisPass,
			db:      env.Envs.RedisDB,
			enabled: env.Envs.RedisEnabled,
		},
		env:    env.Envs.ENV,
		apiURL: env.Envs.ApiUrl,
		mail: mailConfig{
			exp:       env.Envs.MailerExp,
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
				exp:    env.Envs.JwtExp,
				iss:    env.Envs.JwtIss,
			},
		},
		camunda: camundaConfig{
			zeebeAddr:          env.Envs.ZeebeAddr,
			zeebeClientId:      env.Envs.ZeebeClientID,
			zeebeClientSecret:  env.Envs.ZeebeClientSecret,
			zeebeAuthServerUrl: env.Envs.ZeebeAuthServerUrl,
		},
		rateLimiter: ratelimiter.Config{
			RequestPerTimeFrame: env.Envs.RequestPerTimeFrame,
			TimeFrame:           env.Envs.RateLimiterTimeFrame,
			Enabled:             env.Envs.RateLimiterEnabled,
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

	// Cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis cache connection established")
	}

	// Rate Limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	// Mailer
	mailer := mailer.NewSendgrid(cfg.mail.sendgrid.apiKey, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	// zeebe
	zeebeClient, err := zeebe.NewZeebeClient(
		cfg.camunda.zeebeClientId,
		cfg.camunda.zeebeClientSecret,
		cfg.camunda.zeebeAuthServerUrl,
		cfg.camunda.zeebeAddr,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer zeebeClient.Close()

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimiter,
		zeebeClient:   zeebeClient,
	}

	// Metrics Collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutine", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
