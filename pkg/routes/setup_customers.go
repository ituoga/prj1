package routes

import (
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/ituoga/proj1/pkg/experiments"
	"github.com/ituoga/proj1/pkg/pages"
	"github.com/ituoga/proj1/template"
)

func SetupCustomers(router chi.Router, session sessions.Store) {

	var form = pages.CustomerForm{}
	var index = pages.CustomerIndex{}

	router.Route("/customers", func(customersRouter chi.Router) {

		customersRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			index.Data = experiments.Customer().List()
			template.CustomersIndex(index).Render(r.Context(), w)
		})

		customersRouter.Get("/new", func(w http.ResponseWriter, r *http.Request) {
			uid := uuid.NewString()
			form.Form.UUID = &uid
			template.CustomersEdit(form).Render(r.Context(), w)
		})

		customersRouter.Put("/new", func(w http.ResponseWriter, r *http.Request) {
			err := datastar.BodyUnmarshal(r, &form.Form.Data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			sse := datastar.NewSSE(w, r)
			datastar.PatchStore(sse, form.Form.Data)

			if form.Form.UUID == nil {
				uid := uuid.NewString()
				form.Form.UUID = &uid
				experiments.Customer().Store(form.Form.Data, form.Form.UUID)
				return
			}
			experiments.Customer().Update(form.Form)
			datastar.Redirect(sse, "/customers")
		})

		customersRouter.Route("/{uid}", func(customerRouter chi.Router) {
			customerRouter.Get("/edit", func(w http.ResponseWriter, r *http.Request) {
				uid := chi.URLParam(r, "uid")
				form.Form.UUID = &uid
				form.Form = experiments.Customer().Load(*form.Form.UUID)
				template.CustomersEdit(form).Render(r.Context(), w)
			})

			customerRouter.Delete("/delete", func(w http.ResponseWriter, r *http.Request) {
				uid := chi.URLParam(r, "uid")
				experiments.Customer().Delete(uid)
				sse := datastar.NewSSE(w, r)
				datastar.Redirect(sse, "/customers")
			})
		})
	})
}
