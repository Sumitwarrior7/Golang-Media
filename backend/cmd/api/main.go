package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/Sumitwarrior7/social/internal/auth"
	"github.com/Sumitwarrior7/social/internal/db"
	"github.com/Sumitwarrior7/social/internal/env"
	"github.com/Sumitwarrior7/social/internal/mailer"
	"github.com/Sumitwarrior7/social/internal/ratelimiter"
	"github.com/Sumitwarrior7/social/internal/store"
	"github.com/Sumitwarrior7/social/internal/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Social Media Website Documentation
//	@version		1.0
//	@description	This is a sample server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

// @securitydefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		env:         env.GetString("ENV", "development"),
		addr:        env.GetString("ADDR", ":8081"),
		apiUrl:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendUrl: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://sumit_user:pg_pass_key@localhost/social_network?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDDIS_PW", ""),
			db:      env.GetInt("REDDIS_DB", 0),
			enabled: env.GetBool("REDDIS_ENABLED", false),
		},
		mail: mailConfig{
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
			fromEmail: env.GetString("FROM_EMAIL", ""),
			exp:       3 * 24 * time.Hour, // 3 days
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "GolangMedia",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxIdleConns,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleTime,
	)
	// version := 1
	if err != nil {
		logger.Panic(err)
	}
	defer db.Close()
	logger.Info("db connected!!!!!!")

	store := store.NewPostgresStorage(db)

	// Cache
	var redisClient *redis.Client
	if cfg.redisCfg.enabled {
		redisClient = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
	}
	RedisStorage := cache.NewRedisStorage(redisClient)

	// Mailer
	mailerMailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// Authenticator
	JwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	// Rate Limiter
	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailerMailtrap,
		authenticator: JwtAuthenticator,
		cacheStorage:  RedisStorage,
		rateLimiter:   rateLimiter,
	}

	// Metrics/stats to be shown
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("go-routines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
