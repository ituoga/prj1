package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/ituoga/proj1/template"
)

func SetupCustomers(router chi.Router, session sessions.Store) {
	router.Route("/customers", func(customersRouter chi.Router) {
		customersRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			template.CustomersIndex().Render(r.Context(), w)
		})
	})
}
