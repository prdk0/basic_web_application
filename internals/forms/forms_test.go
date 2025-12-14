package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid!")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "b")
	if form.Valid() {
		t.Error("form shows valid when ewquired field missing")
	}
	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")
	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "b")
	if !form.Valid() {
		t.Error("shoes does not have required fields when it does")
	}
}

func TestMinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	postData := url.Values{}
	postData.Add("a", "a")
	r.PostForm = postData
	form := New(r.PostForm)
	form.MinLength("a", 1)
	if !form.Valid() {
		t.Error("form attributes should have the min length of 1")
	}

	isError := form.Errors.Get("a")
	if isError != "" {
		t.Error("should not an error, but got one")
	}

	postData = url.Values{}
	postData.Add("a", "a")
	r.PostForm = postData
	form = New(r.PostForm)
	form.MinLength("a", 2)
	if form.Valid() {
		t.Error("form attributes should have the min length of 2 else fail")
	}
	isError = form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but didnot get one")
	}
}

func TestIsValidEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	postData := url.Values{}
	postData.Add("a", "someone@somewhere.com")
	r.PostForm = postData
	form := New(r.PostForm)
	form.IsValidEmail("a")
	if !form.Valid() {
		t.Error("This test should pass with valid email")
	}

	postData = url.Values{}
	postData.Add("a", "a@")
	r.PostForm = postData
	form = New(r.PostForm)
	form.IsValidEmail("a")
	if form.Valid() {
		t.Error("This test should fail with invalid email")
	}
}

func TestFormHas(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	postData := url.Values{}
	postData.Add("a", "a")
	r.PostForm = postData
	form := New(r.PostForm)
	has := form.Has("a")
	if !has {
		t.Error("Test failed even after field attibutes present")
	}
	postData = url.Values{}
	postData.Add("a", "")
	r.PostForm = postData
	form = New(r.PostForm)
	has = form.Has("a")
	if has {
		t.Error("Test failed even after field attibutes present")
	}
}
