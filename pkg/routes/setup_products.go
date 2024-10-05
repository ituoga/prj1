package routes

import (
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/ituoga/prj1/pkg/experiments"
	"github.com/ituoga/prj1/pkg/pages"
	template "github.com/ituoga/prj1/template"
)

func SetupProducts(router chi.Router, session sessions.Store) {
	var index = pages.ProductsIndex{}
	var form = pages.ProductsForm{}

	router.Route("/products", func(productsRouter chi.Router) {
		productsRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var err error
			index.Data, err = experiments.Product().List()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			template.ProductsIndex(index).Render(r.Context(), w)
		})
		productsRouter.Get("/new", func(w http.ResponseWriter, r *http.Request) {
			uid := uuid.NewString()
			form.Form.UUID = &uid
			form.Form.IsNew = true
			template.ProductsForm(pages.ProductsForm{}).Render(r.Context(), w)
		})

		productsRouter.Put("/store", func(w http.ResponseWriter, r *http.Request) {
			err := datastar.BodyUnmarshal(r, &form.Form.Data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if form.Form.IsNew {
				err = experiments.Product().Create(&form.Form)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else {
				err = experiments.Product().Update(&form.Form)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			sse := datastar.NewSSE(w, r)
			datastar.Redirect(sse, "/products")
		})

		productsRouter.Route("/{id}", func(productRouter chi.Router) {
			productRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				id := chi.URLParam(r, "id")
				var err error
				form.Form, err = experiments.Product().Get(id)
				form.Form.IsNew = false
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				template.ProductsForm(form).Render(r.Context(), w)
			})
		})

	})
}
