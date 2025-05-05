package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/lonewolf101101/Architect-betting/backend/cmd/web/app"
	"github.com/lonewolf101101/Architect-betting/backend/common/oapi"
	"github.com/lonewolf101101/Architect-betting/backend/pkg/customerman"
)

// Create a NoSurf middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly flags set.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // or specific domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle the preflight request (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func IsAuth(r *http.Request) bool {
	isAuth, ok := r.Context().Value(app.ContextKeyIsAuth).(bool)
	if !ok {
		return false
	}
	return isAuth
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsAuth(r) {
			oapi.Forbidden(w)
			return
		}
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.Session.Exists(r, "oauth2_provider_name")
		if !exists || app.Session.GetString(r, "oauth2_provider_name") == "" {
			app.InfoLog.Println("accessToken not found")
			next.ServeHTTP(w, r)
			return
		}

		email := app.Session.GetString(r, "email")
		if len(email) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		customer, err := app.Customers.GetWithEmail(email)
		if err != nil {
			if errors.Is(err, customerman.ErrNotFound) {
				next.ServeHTTP(w, r)
				return
			}
			oapi.ServerError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), app.ContextKeyIsAuth, true)
		ctx = context.WithValue(ctx, app.ContextKeyAuthCustomer, customer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
