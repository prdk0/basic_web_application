package forms

import (
	"net/url"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}), // or make(errors)
	}
}

func (f *Form) Has(fields ...string) bool {
	for _, field := range fields {
		x := f.Get(field)
		if x == "" {
			f.Errors.Add(field, "This field cannot be blank!")
			return false
		}
	}
	return true
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
