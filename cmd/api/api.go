package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/damarteplok/social/docs"
	"github.com/damarteplok/social/internal/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   env.Envs.AllowedOrigin,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(app.RateLimiterMiddleware)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.With(app.BasicAuthMiddleware()).
			Get("/debug/vars", expvar.Handler().ServeHTTP)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})

		// Auth routes
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Get("/", app.getUserAllHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/camunda", func(r chi.Router) {
			r.Route("/resource", func(r chi.Router) {
				r.With(app.BasicAuthMiddleware()).Post("/deploy", app.deployOnlyCamundaHandler)
				r.With(app.BasicAuthMiddleware()).Post("/crud", app.crudCamundaHandler)
				r.With(app.BasicAuthMiddleware()).Post("/deploy-crud", app.deployCamundaHandler)
				r.Post("/{resourceKey}/delete", func(w http.ResponseWriter, r *http.Request) {})
			})
			r.Route("/job", func(r chi.Router) {
				r.Post("/activate", func(w http.ResponseWriter, r *http.Request) {})
				r.Route("/{jobKey}", func(r chi.Router) {
					r.Post("/", func(w http.ResponseWriter, r *http.Request) {})
					r.Patch("/", func(w http.ResponseWriter, r *http.Request) {})
					r.Post("/fail", func(w http.ResponseWriter, r *http.Request) {})
					r.Post("/error", func(w http.ResponseWriter, r *http.Request) {})
				})
			})
			r.Route("/incident", func(r chi.Router) {
				r.Route("/{incidentKey}", func(r chi.Router) {
					r.Post("/", func(w http.ResponseWriter, r *http.Request) {})
				})
			})
			r.Route("/usertask", func(r chi.Router) {
				r.Route("/{usertaskKey}", func(r chi.Router) {
					r.Post("/", func(w http.ResponseWriter, r *http.Request) {})
					r.Patch("/", func(w http.ResponseWriter, r *http.Request) {})
					r.Post("/assignment", func(w http.ResponseWriter, r *http.Request) {})
					r.Post("/unassignment", func(w http.ResponseWriter, r *http.Request) {})
				})
			})
			r.With(app.BasicAuthMiddleware()).Route("/process-instance", func(r chi.Router) {
				r.Post("/", app.createProsesInstance)
				r.Route("/{processinstanceKey}", func(r chi.Router) {
					r.Post("/cancel", app.cancelProcessInstance)
				})
			})
			r.Route("/message", func(r chi.Router) {
				r.Patch("/publish", func(w http.ResponseWriter, r *http.Request) {})
				r.Patch("/correlate", func(w http.ResponseWriter, r *http.Request) {})
			})
			r.Route("/tasklist", func(r chi.Router) {
				r.Post("/", app.searchTaskListHandler)
			})
		})

		r.Route("/bpmn", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Route("/kantor_ngetes_id", func(r chi.Router) {
				r.Post("/", app.createKantorNgetesIdHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("Server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
