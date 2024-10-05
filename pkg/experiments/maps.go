package experiments

import (
	"reflect"
	"strconv"
)

type MapsDB struct{}

func Maps() *MapsDB {
	return &MapsDB{}
}

func inArray(needle string, haystack []Field) int {
	for _, v := range haystack {
		if needle == v.Name {
			return v.Index
		}
	}
	return -1
}

type Field struct {
	Name  string
	Index int
}

func (m MapsDB) ValidFields(p any) []Field {
	typ := reflect.TypeOf(p)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	availFields := make([]Field, 0)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("map")
		availFields = append(availFields, Field{Name: tag, Index: i})
	}

	return availFields
}

func (mm MapsDB) FromMap(m map[string]string, p any) {
	val := reflect.ValueOf(p)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		panic("p must be a pointer to a struct")
	}

	availFields := mm.ValidFields(p)
	for k, v := range m {
		if idx := inArray(k, availFields); idx != -1 {
			fieldValue := val.Field(idx)

			if fieldValue.Kind() == reflect.String {
				fieldValue.SetString(v)
			}

			if fieldValue.Kind() == reflect.Int {
				vv, _ := strconv.Atoi(v)
				fieldValue.SetInt(int64(vv))
			}
			if fieldValue.Kind() == reflect.Float64 {
				vv, _ := strconv.ParseFloat(v, 64)
				fieldValue.SetFloat(vv)
			}
		}
	}
}

func (mm MapsDB) Contains(p any, field string) bool {
	fields := mm.ValidFields(p)
	return inArray(field, fields) != -1
}
