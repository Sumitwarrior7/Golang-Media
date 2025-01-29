package main

import (
	"context"
	"errors"
	"expvar"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/Sumitwarrior7/social/docs"
	"github.com/Sumitwarrior7/social/internal/auth"
	"github.com/Sumitwarrior7/social/internal/mailer"
	"github.com/Sumitwarrior7/social/internal/ratelimiter"
	"github.com/Sumitwarrior7/social/internal/store"
	"github.com/Sumitwarrior7/social/internal/store/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors" // http-swagger middleware
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	cacheStorage  cache.Storage
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiUrl      string
	mail        mailConfig
	frontendUrl string
	auth        authConfig
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
}

/* Redis related configutaions */
type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

/* Authentication related configutaions */
type authConfig struct {
	basic basicConfig
	token tokenConfig
}
type basicConfig struct {
	user string
	pass string
}
type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

/* Email related configutaions */
type mailConfig struct {
	sendGrid  sendGridConfig
	mailTrap  mailTrapConfig
	fromEmail string
	exp       time.Duration
}

// Choice 1
type sendGridConfig struct {
	apiKey string
}

// Choice 2
type mailTrapConfig struct {
	apiKey string
}

/* Database related configutaions */
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// Adding routing to application server
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(app.RateLimiterMiddleware)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Get("/all-users", app.getAllUsersHandler)
		r.With(app.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)

		// docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		// r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.TokenAuthMiddleware())
			r.Post("/", app.createPostHandler)
			r.Get("/", app.getAllPostsHandler)
			r.Get("/user/{userId}", app.getAllPostsByUserIdHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.CheckPostOwnership("admin", app.deletePostHandler))
				r.Patch("/", app.CheckPostOwnership("moderator", app.updatePostHandler))

				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
					r.Route("/{commentId}", func(r chi.Router) {
						r.Use(app.commentContextMiddleware)
						r.Get("/", app.getCommentByIdHandler)
						r.Put("/", app.CheckCommentOwnership("moderator", app.updateCommentHandler))
						r.Delete("/", app.CheckCommentOwnership("admin", app.deleteCommentHandler))
					})
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.TokenAuthMiddleware())
				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.TokenAuthMiddleware())
				r.Get("/feed", app.getUserFeedHandler)
				r.Get("/current-user", app.getCurrentUserHandler)
				r.Get("/followed-users", app.getFollowedUsersHandler)
				// r.Get("/all-users", app.getAllUsersHandler)
			})
		})

		// Public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler) // It is used to create new users
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	// docs.SwaggerInfo.Version = version
	// docs.SwaggerInfo.Host = app.config.apiUrl
	// docs.SwaggerInfo.BasePath = "/v1"

	svr := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	/* Graceful Shutdown */
	// A graceful shutdown in servers ensures that a server completes any ongoing tasks and releases resources properly
	// before stopping.
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("Signal caught", "signal", s.String())

		shutdown <- svr.Shutdown(ctx)
	}()

	app.logger.Infow("Server has started", "addr", app.config.addr, "env", app.config.env)
	err := svr.ListenAndServe() // This blocks the program and keeps the server running, listening for incoming client requests.

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown // Waits for the shutdown channel (from the graceful shutdown goroutine) to send its result.
	if err != nil {
		return err // If an error occurred during the shutdown process (e.g., timeout, resource cleanup failure), it returns the error.
	}

	// server stops finally
	app.logger.Infow("Server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
