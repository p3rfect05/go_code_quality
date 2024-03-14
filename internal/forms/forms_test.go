package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
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
	r := httptest.NewRequest("POST", "/", nil)

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
	r := httptest.NewRequest("POST", "/", nil)

	requiredField := "key1"

	postedData := url.Values{}
	postedData.Add(requiredField, "value1")
	r.PostForm = postedData
	form := New(r.PostForm)
	form.Required(requiredField)
	if !form.Valid() {
		t.Errorf("Field %s should be present in form", requiredField)
	}
	extraField := "lol"
	form.Required(extraField)
	if form.Valid() {
		t.Errorf("Field %s should not be present in form", extraField)
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	postedData := url.Values{}
	invalidEmail := "a@a"
	validEmail := "a@a.com"
	postedData.Add("validEmail", validEmail)
	postedData.Add("invalidEmail", invalidEmail)
	r.PostForm = postedData
	form := New(r.PostForm)
	if form.IsEmail("invalidEmail") {
		t.Errorf("%s should be invalid email", invalidEmail)
	}
	if !form.IsEmail("validEmail") {
		t.Errorf("%s should be valid email", validEmail)
	}
}

func TestForm_MinLength(t *testing.T) {

	r := httptest.NewRequest("POST", "/", nil)

	postedData := url.Values{}
	str := "12345"
	invalidLength := len(str) + 1
	validLength := len(str) - 1

	postedData.Add("str", str)
	r.PostForm = postedData
	form := New(r.PostForm)
	if form.MinLength("str", invalidLength) {
		t.Errorf("%s should have length (minLength = %d)", str, invalidLength)
	}

	if !form.MinLength("str", validLength) {
		t.Errorf("%s should have valid length (minLength = %d)", str, validLength)
	}

}

func TestErrors_Add(t *testing.T) {
	errMap := make(errors)
	startLen := len(errMap)
	errMap.Add("key1", "value1")
	if len(errMap)-startLen != 1 || errMap.Get("key1") == "" {
		t.Errorf("errMap should have one more element key1=value1")
	}
}

func TestErrors_Get(t *testing.T) {
	errMap := make(errors)
	if errMap.Get("key1") != "" {
		t.Errorf("errMap should not contain \"key1\"")
	}
	errMap.Add("key1", "value1")
	if v := errMap.Get("key1"); v != "value1" {
		if len(v) == 0 {
			t.Error("errMap is empty")
		} else {
			t.Error("errMap[\"key1\"] != \"value1\"")
		}
	}
}
