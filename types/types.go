package types

import (
	"context"
	"fmt"
	"log"
	"math"
)

type InvoiceRow struct {
	Number  int
	Name    string
	Price   float64
	Comment string

	UID string

	Quantity float64
	Units    string
	Vat      float64
}

type Invoice struct {
	Number string

	DocumentDate string
	DueDate      string

	Currency   string
	Rate       float64
	SerialName string

	RecipientName    string
	RecipientCode    string
	RecipientVAT     string
	RecipientEmail   string
	RecipientPhone   string
	RecipientAddr    string
	RecipientCountry string

	Lines []InvoiceRow `form:"lines"`

	WrittenBy string
	TakenBy   string
	Note      string
	Comment   string

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
	Title  string
	Signal string
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
	ID     int
	HTMLID string
	Data   T
}
