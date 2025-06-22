package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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

func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	return x != ""
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsValidEmail(field string) bool {
	email := f.Get(field)
	emailRegex := `^[\w\+\.-]+@[a-zA-Z\d\.-]+\.[a-zA-Z]{2,}$`
	rs := regexp.MustCompile(emailRegex)
	if !rs.MatchString(email) {
		f.Errors.Add(field, "Enter a valid email")
		return false
	}
	return true
}
