package routes

import (
	"log"
	"net/http"

	"github.com/delaneyj/datastar"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/ituoga/prj1/pkg/experiments"
	"github.com/ituoga/prj1/pkg/pages"
	template "github.com/ituoga/prj1/template"
)

func SetupAuth(router chi.Router, session sessions.Store) {
	router.Route("/login", func(loginRouter chi.Router) {
		loginRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {

			s, err := session.Get(r, "auth")
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if s.Values["auth"] == true {
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			template.Login(pages.Login{}).Render(r.Context(), w)
		})

		loginRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var pagesLogin pages.Login
			err := datastar.BodyUnmarshal(r, &pagesLogin.DS)
			if err != nil {
				log.Printf("%v", err)
			}

			ok := experiments.Auth(pagesLogin.DS.Username, pagesLogin.DS.Password)

			if ok {
				s, err := session.Get(r, "auth")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				s.Values["auth"] = ok
				s.Values["name"] = pagesLogin.DS.Username
				err = s.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				sse := datastar.NewSSE(w, r)
				datastar.Redirect(sse, "/dashboard")
				return
			}
			datastar.NewSSE(w, r)
		})
	})

	router.Route("/logout", func(logoutRouter chi.Router) {
		logoutRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			s, err := session.Get(r, "auth")
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			s.Values["auth"] = false
			s.Values["name"] = ""
			err = s.Save(r, w)
			if err != nil {
				log.Printf("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		})
	})
}
