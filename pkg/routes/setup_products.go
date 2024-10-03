package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/ituoga/proj1/template"
)

func SetupProducts(router chi.Router, session sessions.Store) {
	router.Route("/products", func(productsRouter chi.Router) {
		productsRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			template.ProductsIndex().Render(r.Context(), w)
		})
	})
}
