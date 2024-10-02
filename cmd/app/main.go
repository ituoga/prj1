package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ituoga/proj1/pkg/routes"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func main() {

	router := chi.NewRouter()

	sessionStore := sessions.NewCookieStore([]byte("secret")) // @todo: move to env
	sessionStore.MaxAge(int(24 * time.Hour / time.Second))

	routes.SetupDefault(router, sessionStore)

	log.Printf("Server started at http://localhost:8089")
	http.ListenAndServe(":8089", router)
}
