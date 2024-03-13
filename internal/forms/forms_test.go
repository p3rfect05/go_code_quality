package forms

import (
	"net/http"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", nil)
	form := New(r.PostForm)
	if !form.Valid() {
		t.Error("Form should be valid")
	}
	form.Has("a")
	if form.Valid() {
		t.Error("Form should be invalid (does not contain \"a\" field)")
	}
}

func TestForm_Has(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", nil)

	form := New(r.PostForm)
	requiredField := "key1"
	if form.Has(requiredField) {
		t.Errorf("Form should not contain %s field", requiredField)
	}
	postedData := url.Values{}
	postedData.Add(requiredField, "value1")
	form = New(postedData)

	if !form.Has(requiredField) {
		t.Errorf("Form should contain %s field", requiredField)
	}
}

func TestForm_Required(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", nil)

	requiredField := "key1"

	postedData := url.Values{}
	postedData.Add(requiredField, "value1")
	form := New(postedData)
	form.Required(requiredField)
	if !form.Valid() {
		t.Errorf("Field %s should be present in form", requiredField)
	}
	form = New(r.PostForm)
	extraField := "lol"
	form.Required(extraField)
	if form.Valid() {
		t.Errorf("Field %s should not be present in form", extraField)
	}
}
