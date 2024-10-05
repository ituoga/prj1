package routes

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ituoga/proj1/pkg/experiments"
	"github.com/ituoga/proj1/types"
)

var settings = types.Settings{
	SerialName: "SR",
}

func init() {
	settings = experiments.Settings().Load(context.TODO())
}

func Middleware(sessionStore sessions.Store) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
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
	}
}
