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
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ := http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	form.Required("d", "e", "f")
	if form.Valid() {
		t.Error("The form shows invalid, but it is valid for sure")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("Should have an error, but did not get one")
	}

	postedValues := url.Values{}
	postedValues.Add("some_field", "some value")
	form = New(postedValues)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows minlength of 100 meet when data is shorter")
	}

	postedValues = url.Values{}
	postedValues.Add("another_field", "abcdefg")
	form = New(postedValues)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("data is longer than what minlength shows")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("Should not have an error, but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("email", "isti@google.com")
	postedData.Add("example", "something")

	form := New(postedData)

	form.IsEmail("another")
	if form.Valid() {
		t.Error("Form shows valid email for non existent ")
	}

	form.IsEmail("email")
	if form.Errors.Get("email") != "" {
		t.Error("Email field validation error, correct email is invalid")
	}

	form.IsEmail("example")
	if form.Errors.Get("example") == "" {
		t.Error("Email field validation error, not correct email adress is valid")
	}
}
