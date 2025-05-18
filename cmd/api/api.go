package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lutif/go-social/internal/config"
	"github.com/lutif/go-social/internal/store"
)

type application struct {
	config config.Config
	store  store.Store
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", app.healthCheckHandler)
	r.Route("/users", func(r chi.Router) {
		r.Post("/", app.createUser)
		r.Get("/feed", app.getFeed)
		r.Route("/{userId}", func(r chi.Router) {
			r.Use(app.userContextMiddleware)
			r.Get("/", app.getUserHandler)
			r.Post("/follow", app.followHandler)
			r.Post("/unfollow", app.unfollowHandler)
		})
	})

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", app.getFeed)
		r.Post("/", app.createPost)
		r.Route("/{postId}",
			func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Post("/", app.updatePost)
				r.Get("/", app.getPost)
				r.Delete("/", app.deletePost)

				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createComment)
					r.Route("/{commentId}",
						func(r chi.Router) {
							r.Use(app.commentContextMiddleware)
							r.Get("/", app.getComment)
							r.Post("/", app.updateComment)
							r.Delete("/", app.deleteComment)
						})

				})
			})

	})
	return r
}

func (app *application) run() error {

	mux := app.mount()

	svr := http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: app.config.WriteTimeout,
		ReadTimeout:  app.config.ReadTimeout,
		IdleTimeout:  app.config.IdleTimeout,
	}
	log.Printf("Sever listing at %s", app.config.Addr)
	return svr.ListenAndServe()
}
