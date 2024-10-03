package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ituoga/proj1/template"
	"github.com/ituoga/proj1/types"

	"github.com/delaneyj/datastar"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nuid"
)

func SetupDefault(sr chi.Router, session sessions.Store) {

	var i = 10

	var invoiceForm = types.Invoice{
		Rate:         1.0,
		Currency:     "EUR",
		SerialName:   "SR",
		DocumentDate: time.Now().Format("2006-01-02"),
		DueDate:      time.Now().AddDate(0, 0, 15).Format("2006-01-02"),
		Lines:        []types.InvoiceRow{{UID: nuid.Next(), Units: "vnt", Vat: 1, Quantity: 1}},
	}
	var settings = types.Settings{
		SerialName: "SR",
	}

	var autocomplete = []types.Complete{
		{Signal: "1", Title: "Surname 1"},
		{Signal: "2", Title: "Name 2"},
		{Signal: "3", Title: "Name 3"},
		{Signal: "4", Title: "Name 4"},
		{Signal: "5", Title: "Name 5"},
		{Signal: "6", Title: "Name 6"},
	}

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
		invoiceForm.Lines = append(invoiceForm.Lines, types.InvoiceRow{UID: nuid.Next(), Units: "vnt", Vat: 1, Quantity: 1})

		sse := datastar.NewSSE(w, r)

		datastar.RenderFragmentTempl(sse, template.InvoiceLine(
			fmt.Sprintf("%d", len(invoiceForm.Lines)), invoiceForm.Lines[len(invoiceForm.Lines)-1],
		), datastar.WithoutViewTransitions(),
			datastar.WithMergeAppendElement(),
			datastar.WithQuerySelector("tbody"),
		)
		datastar.RenderFragmentTempl(sse, template.InvoiceLine2(
			fmt.Sprintf("%d", len(invoiceForm.Lines)), invoiceForm.Lines[len(invoiceForm.Lines)-1],
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
		for i, v := range invoiceForm.Lines {
			if v.UID != mpp["eln"].(string) {
				lines = append(lines, invoiceForm.Lines[i])
			}
		}

		invoiceForm.Lines = lines

		invoiceForm.VAT()

		datastar.Delete(sse, fmt.Sprintf("#row-%s-one", id), datastar.WithoutViewTransitions())
		datastar.Delete(sse, fmt.Sprintf("#row-%s-two", id), datastar.WithoutViewTransitions())
		datastar.RenderFragmentTempl(sse, template.Summary(invoiceForm.Summary), datastar.WithoutViewTransitions())
		// datastar.PatchStore(sse, mpp)
	})

	sr.HandleFunc("/complete/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		_ = id
		sse := datastar.NewSSE(w, r)
		_ = sse

		for _, v := range autocomplete {
			if v.Signal == id {
				invoiceForm.RecipientName = v.Title
			}
		}

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

		for lno, line := range invoiceForm.Lines {
			if line.UID == mpp.Eln {
				switch mpp.Elf {
				case "name":
					invoiceForm.Lines[lno].Name = mpp.Elv
				case "price":
					invoiceForm.Lines[lno].Price, _ = strconv.ParseFloat(mpp.Elv, 64)
				case "comment":
					invoiceForm.Lines[lno].Comment = mpp.Elv
				case "qty":
					invoiceForm.Lines[lno].Quantity, _ = strconv.ParseFloat(mpp.Elv, 64)
				case "units":
					invoiceForm.Lines[lno].Units = mpp.Elv
				case "vat":
					invoiceForm.Lines[lno].Vat, _ = strconv.ParseFloat(mpp.Elv, 64)
				}
				tdid := fmt.Sprintf("%s", mpp.Eln)
				// log.Printf("%v", tdid)
				invoiceForm.VAT()
				vat := invoiceForm.Lines[lno].Price * invoiceForm.Lines[lno].Quantity * invoiceForm.Lines[lno].Vat
				_ = vat
				datastar.RenderFragmentTempl(sse, template.TDAmount(tdid, vat), datastar.WithoutViewTransitions(), datastar.WithQuerySelectorID("td-id-"+mpp.Eln))
				// datastar.RenderFragmentTempl(sse, template.Summary(invoiceForm.Summary), datastar.WithoutViewTransitions())
				// datastar.RenderFragmentTempl(sse, template.Debug(invoiceForm), datastar.WithoutViewTransitions())
			}
		}

		switch mpp.Elf {
		case "document_date":
			invoiceForm.DocumentDate = mpp.Elv
		case "due_date":
			invoiceForm.DueDate = mpp.Elv
		case "currency":
			invoiceForm.Currency = mpp.Elv
		case "currency_rate":
			invoiceForm.Rate, _ = strconv.ParseFloat(mpp.Elv, 64)
		case "serial_name":
			invoiceForm.SerialName = mpp.Elv
		case "recipient_name":
			invoiceForm.RecipientName = mpp.Elv
		case "recipient_code":
			invoiceForm.RecipientCode = mpp.Elv
		case "recipient_vat":
			invoiceForm.RecipientVAT = mpp.Elv
		case "recipient_email":
			invoiceForm.RecipientEmail = mpp.Elv
		case "recipient_phone":
			invoiceForm.RecipientPhone = mpp.Elv
		case "recipient_addr":
			invoiceForm.RecipientAddr = mpp.Elv
		case "recipient_country":
			invoiceForm.RecipientCountry = mpp.Elv
		case "written_by":
			invoiceForm.WrittenBy = mpp.Elv
		case "taken_by":
			invoiceForm.TakenBy = mpp.Elv
		case "note":
			invoiceForm.Note = mpp.Elv
		}

		tmp := []types.Complete{}
		if invoiceForm.RecipientName != "" {
			for _, v := range autocomplete {
				if strings.Contains(strings.ToLower(v.Title), strings.ToLower(invoiceForm.RecipientName)) {
					if len(v.Title) == len(invoiceForm.RecipientName) {
						continue
					}
					tmp = append(tmp, v)
				}
			}
		}

		datastar.RenderFragmentTempl(sse, template.Autocomplete(tmp, "pirkejas"), datastar.WithoutViewTransitions())

		datastar.RenderFragmentTempl(sse, template.Summary(invoiceForm.Summary), datastar.WithoutViewTransitions())

		datastar.RenderFragmentTempl(sse, template.Debug(invoiceForm), datastar.WithoutViewTransitions())

		_ = sse
	})

	sr.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {

			var mpp map[string]interface{}
			err := datastar.BodyUnmarshal(r, &mpp)
			if err != nil {
				log.Printf("%v", err)
				return
			}
			switch mpp["elf"] {
			case "my-name":
				settings.MyName = mpp["elv"].(string)
			case "my-code":
				settings.MyCode = mpp["elv"].(string)
			case "my-vat":
				settings.MyVAT = mpp["elv"].(string)
			case "my-email":
				settings.MyEmail = mpp["elv"].(string)
			case "my-phone":
				settings.MyPhone = mpp["elv"].(string)
			case "my-address":
				settings.MyAddr = mpp["elv"].(string)
			case "my-country":
				settings.MyCountry = mpp["elv"].(string)
			case "my-series":
				settings.SerialName = mpp["elv"].(string)
			case "my-number":
				settings.SerialNo, _ = strconv.Atoi(mpp["elv"].(string))
			}

			datastar.NewSSE(w, r)
			return
		}

		template.Settings(settings, map[string]interface{}{
			"elv": "", "eln": "", "elf": "",
		}).Render(r.Context(), w)
	})

	sr.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		invoiceForm = types.Invoice{
			Rate:         1.0,
			Currency:     "EUR",
			SerialName:   "SR",
			DocumentDate: time.Now().Format("2006-01-02"),
			DueDate:      time.Now().AddDate(0, 0, 15).Format("2006-01-02"),
			Lines:        []types.InvoiceRow{{UID: nuid.Next(), Units: "vnt", Vat: 1, Quantity: 1}},
		}
		invoiceForm.VAT()

		sse := datastar.NewSSE(w, r)
		datastar.RenderFragmentTempl(sse, template.InvoiceLines(invoiceForm.Lines, nil), datastar.WithoutViewTransitions())
		datastar.RenderFragmentTempl(sse, template.Summary(invoiceForm.Summary), datastar.WithoutViewTransitions())

	})

	sr.HandleFunc("/form", func(w http.ResponseWriter, r *http.Request) {
		invoiceForm.VAT()
		// log.Printf("%+v", invoiceForm)
		template.Index2(fmt.Sprintf("%d", i), map[string]interface{}{
			"elv": "", "eln": "", "elf": "",
		}, invoiceForm).Render(r.Context(), w)
	})

	sr.Route("/invoices", func(invoicesRouter chi.Router) {
		invoicesRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			template.InvoicesIndex().Render(r.Context(), w)
		})
	})
}
