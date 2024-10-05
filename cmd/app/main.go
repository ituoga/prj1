package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ituoga/prj1/pkg/routes"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func main() {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	sessionStore := sessions.NewCookieStore([]byte("secret")) // @todo: move to env
	sessionStore.MaxAge(int(24 * time.Hour / time.Second))

	routes.SetupAuth(router, sessionStore)

	authRouter := router.Group(func(r chi.Router) {
		r.Use(routes.Middleware(sessionStore))
	})

	routes.SetupDefault(authRouter, sessionStore)
	routes.SetupCustomers(authRouter, sessionStore)
	routes.SetupProducts(authRouter, sessionStore)
	routes.SetupSettings(authRouter, sessionStore)

	log.Printf("Server started at http://localhost:8089")
	http.ListenAndServe(":8089", router)
}
