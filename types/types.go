package types

import (
	"context"
	"fmt"
	"log"
	"math"
)

type InvoiceRow struct {
	Number  int     `map:"number"`
	Name    string  `map:"name"`
	Price   float64 `map:"price"`
	Comment string  `map:"comment"`

	UID string

	Quantity float64 `map:"qty"`
	Units    string  `map:"units"`
	Vat      float64 `map:"vat"`
}

type Invoice struct {
	Number string `map:"number"`

	DocumentDate string `map:"document_date"`
	DueDate      string `map:"due_date"`

	Currency   string  `map:"currency"`
	Rate       float64 `map:"currency_rate"`
	SerialName string  `map:"serial_name"`

	RecipientName    string `map:"recipient_name"`
	RecipientCode    string `map:"recipient_code"`
	RecipientVAT     string `map:"recipient_vat"`
	RecipientEmail   string `map:"recipient_email"`
	RecipientPhone   string `map:"recipient_phone"`
	RecipientAddr    string `map:"recipient_addr"`
	RecipientCountry string `map:"recipient_country"`

	Lines []InvoiceRow

	WrittenBy string `map:"written_by"`
	TakenBy   string `map:"taken_by"`
	Note      string `map:"note"`
	Comment   string `map:"comment"`

	Summary Summary
}

func (i *Invoice) VAT() {
	vat := 0.0
	total := 0.0
	pretotal := 0.0
	i.Summary = Summary{}

	for _, line := range i.Lines {
		var tmpTotal = line.Price * line.Quantity
		var tmpVat = tmpTotal*line.Vat - tmpTotal

		vat += tmpVat
		pretotal += tmpTotal
		total += tmpTotal + tmpVat

		i.Summary.Total = total
		i.Summary.Subtotal = pretotal
		i.Summary.VATTotal = vat

		i.Summary.AddVat(line.Vat, tmpVat)
	}
}

type Vat struct {
	Name  string
	Value float64
}

type Summary struct {
	Total    float64
	Subtotal float64
	Vat      []Vat
	VATTotal float64
}

type Complete struct {
	Title  string `json:"title" db:"name"`
	Signal string `json:"signal" db:"id"`
}

func (s *Summary) AddVat(in float64, value float64) {
	name := fmt.Sprintf("%d", int(math.Round(in*100)-100))
	if in == 1 {
		return
	}
	for i, v := range s.Vat {
		if v.Name == name {
			s.Vat[i].Value += value
			return
		}
	}
	s.Vat = append(s.Vat, Vat{name, value})
}

type Translations struct {
	Data map[string]string
}

func (t *Translations) Get(key string) string {
	log.Printf("%v", t)
	val, ok := t.Data[key]
	if !ok {
		return key
	}
	return val
}

func Translation(ctx context.Context) *Translations {
	t, ok := ctx.Value("translations").(*Translations)
	if !ok {
		return &Translations{}
	}
	return t

}

type DataStore struct {
	Elf string `json:"elf"`
	Elv string `json:"elv"`
	Eln string `json:"eln"`
}

type Settings struct {
	MyName    string
	MyCode    string
	MyVAT     string
	MyEmail   string
	MyPhone   string
	MyAddr    string
	MyCountry string

	PersonalActivity bool

	SerialName string
	SerialNo   int
}

type Store[T any] struct {
	ID     int     `json:"id" db:"id"`
	UUID   *string `json:"uuid" db:"uuid"`
	HTMLID *string
	IsNew  bool
	Data   T `json:"data"`
}

type Customer struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	VAT     string `json:"vat"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Addr    string `json:"addr"`
	Country string `json:"country"`
}

type CustomerStore Store[Customer]

type Product struct {
	Name string `json:"name" db:"name"`
	Code string `json:"code" db:"code"`
}

type ProductStore Store[Product]

func (p *ProductStore) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":   p.ID,
		"uuid": p.UUID,
		"name": p.Data.Name,
		"code": p.Data.Code,
	}
}
