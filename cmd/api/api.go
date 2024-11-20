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
		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/health", app.healthCheckHandler)
		})

		r.With(app.BasicAuthMiddleware()).
			Get("/debug/vars", expvar.Handler().ServeHTTP)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		// Public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/user", app.getTokenUserHandler)
			})
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
			r.Use(app.AuthTokenMiddleware)
			r.Route("/resource", func(r chi.Router) {
				r.Post("/deploy", app.deployOnlyCamundaHandler)
				r.Post("/crud", app.crudCamundaHandler)
				r.Post("/deploy-crud", app.deployCamundaHandler)
				r.Post("/{processDefinitionKey}/delete", app.deleteCamundaHandler)
				r.Get("/{processDefinitionKey}/xml", app.xmlCamundaHandler)
				r.Get("/operate/statistics", app.operateStatisticsHandler)
			})
			r.Route("/incident", func(r chi.Router) {
				r.Route("/{incidentKey}", func(r chi.Router) {
					r.Post("/resolve", app.resolveIncidentHandler)
				})
			})
			r.Route("/process-instance", func(r chi.Router) {
				r.Post("/", app.createProsesInstance)
				r.Get("/", app.searchProcessInstance)
				r.Route("/{processinstanceKey}", func(r chi.Router) {
					r.Post("/cancel", app.cancelProcessInstance)
				})
			})
			r.Route("/user-task", func(r chi.Router) {
				r.Post("/", app.searchTaskListHandler)
				r.Post("/search", app.searchUserTaskHandler)
			})
		})

		r.Route("/bpmn", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			// GENERATE ROUTES API

			r.Route("/pembuatan_media_berita_technology", func(r chi.Router) {
				r.Get("/", app.searchPembuatanMediaBeritaTechnologyHandler)
				r.Post("/", app.createPembuatanMediaBeritaTechnologyHandler)
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", app.getByIdPembuatanMediaBeritaTechnologyHandler)
					r.Delete("/", app.cancelPembuatanMediaBeritaTechnologyHandler)
					r.Patch("/", app.updatePembuatanMediaBeritaTechnologyHandler)
					r.Get("/history", app.getHistoryByIdPembuatanMediaBeritaTechnologyHandler)
					r.Get("/incidents", app.getIncidentsByIdPembuatanMediaBeritaTechnologyHandler)
				})
			})

			// GENERATE USER TASK ROUTES API

			r.Route("/approvingartikel", func(r chi.Router) {
			})

			r.Route("/reviewingartikel", func(r chi.Router) {
			})

			r.Route("/pembuatanartikel", func(r chi.Router) {
				r.Get("/", app.getUserTaskActivePembuatanArtikelHandler)
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
