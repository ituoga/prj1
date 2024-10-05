package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/ituoga/prj1/pkg/experiments"
	"github.com/ituoga/prj1/pkg/pages"
	template "github.com/ituoga/prj1/template"
	types "github.com/ituoga/prj1/types"

	"github.com/delaneyj/datastar"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nuid"
)

func SetupDefault(sr chi.Router, session sessions.Store) {

	var i = 10

	SetField := func(obj interface{}, name string, value interface{}) error {
		structValue := reflect.ValueOf(obj).Elem()
		structFieldValue := structValue.FieldByName(name)

		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in obj", name)
		}

		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", name)
		}

		structFieldType := structFieldValue.Type()
		val := reflect.ValueOf(value)
		if structFieldType != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
		}

		structFieldValue.Set(val)
		return nil
	}
	_ = SetField

	var invoiceForm = pages.InvoiceForm{
		Form: types.Invoice{
			Rate:         1.0,
			Currency:     "EUR",
			SerialName:   "SR",
			DocumentDate: time.Now().Format("2006-01-02"),
			DueDate:      time.Now().AddDate(0, 0, 15).Format("2006-01-02"),
			Lines:        []types.InvoiceRow{{UID: nuid.Next(), Units: "vnt", Vat: 1, Quantity: 1}},
		},
		Settings: &settings,
	}

	var autocomplete = []types.Complete{
		{Signal: "1", Title: "Surname 1"},
		{Signal: "2", Title: "Name 2"},
		{Signal: "3", Title: "Name 3"},
		{Signal: "4", Title: "Name 4"},
		{Signal: "5", Title: "Name 5"},
		{Signal: "6", Title: "Name 6"},
	}
	_ = autocomplete

	sr.Route("/dashboard", func(dashboardRouter chi.Router) {
		dashboardRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/invoices", http.StatusSeeOther)
		})
	})

	n := 0
	sr.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		t := time.NewTicker(1 * time.Second)

		t2 := time.NewTimer(3 * time.Second)

		for {
			select {
			case <-r.Context().Done():
				log.Printf("client closed connection")
				return
			case <-t2.C:
				t2.Reset(3 * time.Second)
				i++
				datastar.RenderFragmentTempl(sse, template.Vienas(fmt.Sprintf("%d", i)), datastar.WithoutViewTransitions())
			case <-t.C:
				t.Reset(1 * time.Second)
				n++
				datastar.RenderFragmentTempl(sse, template.Du(fmt.Sprintf("%d", n)), datastar.WithoutViewTransitions())
			}
		}
	})

	sr.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		i++
		sse := datastar.NewSSE(w, r)
		datastar.RenderFragmentTempl(sse, template.Vienas(fmt.Sprintf("%d", i)), datastar.WithoutViewTransitions())
		datastar.RenderFragmentTempl(sse, template.Du(fmt.Sprintf("%d", i)), datastar.WithoutViewTransitions())
	})

	sr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// template.Index(fmt.Sprintf("%d", i)).Render(context.TODO(), w)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})

	sr.HandleFunc("/signal", func(w http.ResponseWriter, r *http.Request) {
		invoiceForm.Form.Lines = append(invoiceForm.Form.Lines, types.InvoiceRow{UID: nuid.Next(), Units: "vnt", Vat: 1, Quantity: 1})

		sse := datastar.NewSSE(w, r)

		datastar.RenderFragmentTempl(sse, template.InvoiceLine(
			fmt.Sprintf("%d", len(invoiceForm.Form.Lines)), invoiceForm.Form.Lines[len(invoiceForm.Form.Lines)-1],
		), datastar.WithoutViewTransitions(),
			datastar.WithMergeAppendElement(),
			datastar.WithQuerySelector("tbody"),
		)
		datastar.RenderFragmentTempl(sse, template.InvoiceLine2(
			fmt.Sprintf("%d", len(invoiceForm.Form.Lines)), invoiceForm.Form.Lines[len(invoiceForm.Form.Lines)-1],
		), datastar.WithoutViewTransitions(),
			datastar.WithMergeAppendElement(),
			datastar.WithQuerySelector("#tbody"),
		)
		_ = sse

	})

	sr.HandleFunc("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		log.Printf("id : %v", id)

		var mpp = map[string]interface{}{}
		datastar.BodyUnmarshal(r, &mpp)

		sse := datastar.NewSSE(w, r)

		var lines []types.InvoiceRow
		for i, v := range invoiceForm.Form.Lines {
			if v.UID != mpp["eln"].(string) {
				lines = append(lines, invoiceForm.Form.Lines[i])
			}
		}

		invoiceForm.Form.Lines = lines

		invoiceForm.Form.VAT()

		datastar.Delete(sse, fmt.Sprintf("#row-%s-one", id), datastar.WithoutViewTransitions())
		datastar.Delete(sse, fmt.Sprintf("#row-%s-two", id), datastar.WithoutViewTransitions())
		datastar.RenderFragmentTempl(sse, template.Summary(invoiceForm.Form.Summary), datastar.WithoutViewTransitions())
		// datastar.PatchStore(sse, mpp)
	})

	sr.HandleFunc("/complete/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		_ = id
		sse := datastar.NewSSE(w, r)
		_ = sse

		autocomplete, err := experiments.AutoComplete().Get(r.Context(), id)
		if err != nil {
			log.Printf("?? %v", err)
		}

		invoiceForm.Form.RecipientName = autocomplete.Data.Name
		invoiceForm.Form.RecipientCode = autocomplete.Data.Code
		invoiceForm.Form.RecipientVAT = autocomplete.Data.VAT
		invoiceForm.Form.RecipientAddr = autocomplete.Data.Addr
		invoiceForm.Form.RecipientCountry = autocomplete.Data.Country
		invoiceForm.Form.RecipientPhone = autocomplete.Data.Phone
		invoiceForm.Form.RecipientEmail = autocomplete.Data.Email

		datastar.RenderFragmentTempl(sse, template.InvoiceRecipient(invoiceForm), datastar.WithoutViewTransitions())

		datastar.RenderFragmentTempl(sse, template.Debug(invoiceForm), datastar.WithoutViewTransitions())
	})

	sr.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {

		var mpp types.DataStore
		err := datastar.BodyUnmarshal(r, &mpp)
		sse := datastar.NewSSE(w, r)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		for lno, line := range invoiceForm.Form.Lines {
			if line.UID == mpp.Eln {
				experiments.Maps().FromMap(map[string]string{mpp.Elf: mpp.Elv}, &invoiceForm.Form.Lines[lno])

				tdid := fmt.Sprintf("%s", mpp.Eln)

				invoiceForm.Form.VAT()
				vat := invoiceForm.Form.Lines[lno].Price * invoiceForm.Form.Lines[lno].Quantity * invoiceForm.Form.Lines[lno].Vat
				_ = vat
				datastar.RenderFragmentTempl(sse, template.TDAmount(tdid, vat), datastar.WithoutViewTransitions(), datastar.WithQuerySelectorID("td-id-"+mpp.Eln))
				// datastar.RenderFragmentTempl(sse, template.Summary(invoiceForm.Form.Summary), datastar.WithoutViewTransitions())
				// datastar.RenderFragmentTempl(sse, template.Debug(invoiceForm), datastar.WithoutViewTransitions())
			}
		}

		experiments.Maps().FromMap(map[string]string{mpp.Elf: mpp.Elv}, &invoiceForm.Form)

		tmp := []types.Complete{}
		tmp, err = experiments.AutoComplete().List(invoiceForm.Form.RecipientName)
		if err != nil {
			log.Printf("??? %v", err)
		}
		log.Printf("ac1 %v %v", tmp, invoiceForm.Form.RecipientName)
		// if invoiceForm.Form.RecipientName != "" {
		// 	for _, v := range autocomplete {
		// 		if strings.Contains(strings.ToLower(v.Title), strings.ToLower(invoiceForm.Form.RecipientName)) {
		// 			if len(v.Title) == len(invoiceForm.Form.RecipientName) {
		// 				continue
		// 			}
		// 			tmp = append(tmp, v)
		// 		}
		// 	}
		// }

		datastar.RenderFragmentTempl(sse, template.Autocomplete(tmp, "pirkejas"), datastar.WithoutViewTransitions())

		datastar.RenderFragmentTempl(sse, template.Summary(invoiceForm.Form.Summary), datastar.WithoutViewTransitions())

		datastar.RenderFragmentTempl(sse, template.Debug(invoiceForm), datastar.WithoutViewTransitions())

		_ = sse
	})

	sr.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		invoiceForm.Form = types.Invoice{
			Rate:         1.0,
			Currency:     "EUR",
			SerialName:   "SR",
			DocumentDate: time.Now().Format("2006-01-02"),
			DueDate:      time.Now().AddDate(0, 0, 15).Format("2006-01-02"),
			Lines:        []types.InvoiceRow{{UID: nuid.Next(), Units: "vnt", Vat: 1, Quantity: 1}},
		}

		invoiceForm.Form.VAT()

		sse := datastar.NewSSE(w, r)
		datastar.RenderFragmentTempl(sse, template.InvoiceLines(invoiceForm.Form.Lines, nil), datastar.WithoutViewTransitions())
		datastar.RenderFragmentTempl(sse, template.Summary(invoiceForm.Form.Summary), datastar.WithoutViewTransitions())

	})

	sr.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		invoiceForm.Form.VAT()

		template.InvoiceForm(invoiceForm).Render(r.Context(), w)
	})

	sr.Put("/store", func(w http.ResponseWriter, r *http.Request) {
		experiments.Invoice().Store(r.Context(), invoiceForm.Form)
		sse := datastar.NewSSE(w, r)
		_ = sse
	})

	sr.Route("/invoices", func(invoicesRouter chi.Router) {
		invoicesRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			invoices, _ := experiments.Invoice().List(r.Context())
			template.InvIndex(invoices).Render(r.Context(), w)
		})

		invoicesRouter.Get("/new", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("new page"))
		})

		invoicesRouter.Route("/{id}", func(invoiceRouter chi.Router) {
			getId := func(r *http.Request) string {
				id := r.PathValue("id")
				return id
			}
			invoiceRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				id := getId(r)
				w.Write([]byte(fmt.Sprintf("invoice id: %s", id)))
			})
			invoiceRouter.Get("/edit", func(w http.ResponseWriter, r *http.Request) {
				id := getId(r)
				data, err := experiments.Invoice().Load(r.Context(), id)
				if err != nil {
					log.Printf("%v", err)
				}

				invoiceForm.Form = data.Data
				if len(invoiceForm.Form.Lines) == 0 {
					invoiceForm.Form.Lines = []types.InvoiceRow{{UID: nuid.Next(), Units: "vnt", Vat: 1.00, Quantity: 1}}
				}

				invoiceForm.Form.VAT()

				template.InvoiceForm(invoiceForm).Render(r.Context(), w)
			})
		})
	})

}
