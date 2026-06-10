package forms

import (
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	data := url.Values{}
	f := New(data)
	if f == nil {
		t.Fatal("expected non-nil Form")
	}
	if f.Errors == nil {
		t.Fatal("expected non-nil Errors map")
	}
	if len(f.Errors) != 0 {
		t.Fatalf("expected empty errors, got %d", len(f.Errors))
	}
}

func TestNewWithValues(t *testing.T) {
	data := url.Values{}
	data.Set("username", "john")
	f := New(data)
	if f.Get("username") != "john" {
		t.Fatalf("expected 'john', got '%s'", f.Get("username"))
	}
}

func TestRequired_AddsErrorOnBlank(t *testing.T) {
	f := New(url.Values{})
	f.Set("username", "")
	f.Set("password", "")
	f.Required("username", "password")

	if f.Valid() {
		t.Fatal("expected invalid when required fields are blank")
	}
	if f.Errors.Get("username") == "" {
		t.Fatal("expected error for username")
	}
	if f.Errors.Get("password") == "" {
		t.Fatal("expected error for password")
	}
}

func TestRequired_NoErrorOnNonBlank(t *testing.T) {
	f := New(url.Values{})
	f.Set("username", "john")
	f.Set("password", "secret")
	f.Required("username", "password")

	if !f.Valid() {
		t.Fatal("expected valid when required fields are filled")
	}
}

func TestValid_TrueWithNoErrors(t *testing.T) {
	f := New(url.Values{})
	if !f.Valid() {
		t.Fatal("expected valid for fresh form")
	}
}

func TestValid_FalseWithErrors(t *testing.T) {
	f := New(url.Values{})
	f.Errors.Add("field", "error")
	if f.Valid() {
		t.Fatal("expected invalid when errors exist")
	}
}

func TestMaxLength_NoErrorOnEmpty(t *testing.T) {
	f := New(url.Values{})
	f.MaxLength("field", 5)
	if !f.Valid() {
		t.Fatal("expected valid when field is empty")
	}
}

func TestMaxLength_NoErrorWithinLimit(t *testing.T) {
	f := New(url.Values{})
	f.Set("field", "abc")
	f.MaxLength("field", 5)
	if !f.Valid() {
		t.Fatal("expected valid when within max length")
	}
}

func TestMaxLength_ErrorWhenExceeded(t *testing.T) {
	f := New(url.Values{})
	f.Set("field", "abcdef")
	f.MaxLength("field", 5)
	if f.Valid() {
		t.Fatal("expected invalid when exceeds max length")
	}
}

func TestMinLength_NoErrorOnEmpty(t *testing.T) {
	f := New(url.Values{})
	f.Minlength("field", 5)
	if !f.Valid() {
		t.Fatal("expected valid when field is empty")
	}
}

func TestMinLength_NoErrorWithinLimit(t *testing.T) {
	f := New(url.Values{})
	f.Set("field", "abcdef")
	f.Minlength("field", 5)
	if !f.Valid() {
		t.Fatal("expected valid when within min length")
	}
}

func TestMinLength_ErrorWhenTooShort(t *testing.T) {
	f := New(url.Values{})
	f.Set("field", "ab")
	f.Minlength("field", 5)
	if f.Valid() {
		t.Fatal("expected invalid when below min length")
	}
}

func TestPermittedValues_NoErrorOnEmpty(t *testing.T) {
	f := New(url.Values{})
	f.PermittedValues("field", "a", "b", "c")
	if !f.Valid() {
		t.Fatal("expected valid when field is empty")
	}
}

func TestPermittedValues_NoErrorOnValid(t *testing.T) {
	f := New(url.Values{})
	f.Set("field", "b")
	f.PermittedValues("field", "a", "b", "c")
	if !f.Valid() {
		t.Fatal("expected valid when value is permitted")
	}
}

func TestPermittedValues_ErrorOnInvalid(t *testing.T) {
	f := New(url.Values{})
	f.Set("field", "d")
	f.PermittedValues("field", "a", "b", "c")
	if f.Valid() {
		t.Fatal("expected invalid when value is not permitted")
	}
}

func TestMatchesPattern_NoErrorOnEmpty(t *testing.T) {
	f := New(url.Values{})
	f.MatchesPattern("field", EmailRX)
	if !f.Valid() {
		t.Fatal("expected valid when field is empty")
	}
}

func TestMatchesPattern_NoErrorOnMatch(t *testing.T) {
	f := New(url.Values{})
	f.Set("field", "user@example.com")
	f.MatchesPattern("field", EmailRX)
	if !f.Valid() {
		t.Fatal("expected valid when value matches pattern")
	}
}

func TestMatchesPattern_ErrorOnNoMatch(t *testing.T) {
	f := New(url.Values{})
	f.Set("field", "not-an-email")
	f.MatchesPattern("field", EmailRX)
	if f.Valid() {
		t.Fatal("expected invalid when value does not match pattern")
	}
}
