package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/delaneyj/datastar"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/ituoga/prj1/pkg/experiments"
	template "github.com/ituoga/prj1/template"
	types "github.com/ituoga/prj1/types"
)

func SetupSettings(sr chi.Router, session sessions.Store) {

	sr.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {

			var mpp types.DataStore
			err := datastar.BodyUnmarshal(r, &mpp)
			if err != nil {
				log.Printf("%v", err)
				return
			}
			switch mpp.Elf {
			case "my-name":
				settings.MyName = mpp.Elv
			case "my-code":
				settings.MyCode = mpp.Elv
			case "my-vat":
				settings.MyVAT = mpp.Elv
			case "my-email":
				settings.MyEmail = mpp.Elv
			case "my-phone":
				settings.MyPhone = mpp.Elv
			case "my-address":
				settings.MyAddr = mpp.Elv
			case "my-country":
				settings.MyCountry = mpp.Elv
			case "my-series":
				settings.SerialName = mpp.Elv
			case "my-number":
				settings.SerialNo, _ = strconv.Atoi(mpp.Elv)
			}

			datastar.NewSSE(w, r)
			return
		}
		if r.Method == http.MethodPut {
			sse := datastar.NewSSE(w, r)
			_ = sse
			experiments.Settings().Store(r.Context(), settings)
			return
		}

		template.Settings(settings, map[string]interface{}{
			"elv": "", "eln": "", "elf": "",
		}).Render(r.Context(), w)
	})
}
