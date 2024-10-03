package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ituoga/proj1/pkg/experiments"
	"github.com/ituoga/proj1/pkg/routes"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func main() {

	experiments.DB()

	router := chi.NewRouter()

	sessionStore := sessions.NewCookieStore([]byte("secret")) // @todo: move to env
	sessionStore.MaxAge(int(24 * time.Hour / time.Second))

	routes.SetupAuth(router, sessionStore)
	authRouter := router.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				session, err := sessionStore.Get(r, "auth")
				if err != nil {
					log.Printf("%v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return

				}

				if session.Values["auth"] != true {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}

				ctx := context.WithValue(r.Context(), "user", session.Values["name"].(string))
				h.ServeHTTP(w, r.WithContext(ctx))
			})
		})
	})

	routes.SetupDefault(authRouter, sessionStore)
	routes.SetupCustomers(authRouter, sessionStore)
	routes.SetupProducts(authRouter, sessionStore)

	log.Printf("Server started at http://localhost:8089")
	http.ListenAndServe(":8089", router)
}
