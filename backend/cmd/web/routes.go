package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lonewolf101101/Architect-betting/backend/cmd/web/app"
)

func routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Use(SecureHeaders, app.Session.Enable, Authenticate, corsMiddleware)

	// if app.Mode == "debug" {
	// 	cr.Get("/api/ws", app.FrontendWS.Handler)
	// } else {
	// r.With(RequireAuth).Get("/api/ws", app.FrontendWS.Handler)
	// }
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/pub", func(r chi.Router) {
		r.Get("/visit", customerVisit)
		r.Route("/auth", func(r chi.Router) {
			r.Get("/login", oauthLogin(app.Google))
			r.Get("/callback", oauthCallback(app.Google))
		})
	})

	r.With(RequireAuth).Route("/api", func(r chi.Router) {
		r.Get("/me", Me)
		r.Get("/logout", logout)

	})

	return r
}
