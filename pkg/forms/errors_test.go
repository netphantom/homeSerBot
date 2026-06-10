package forms

import (
	"testing"
)

func TestErrorsAdd(t *testing.T) {
	e := make(errors)
	e.Add("field1", "error message 1")
	e.Add("field1", "error message 2")
	e.Add("field2", "error message 3")

	if len(e) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(e))
	}
	if len(e["field1"]) != 2 {
		t.Fatalf("expected 2 errors for field1, got %d", len(e["field1"]))
	}
	if e["field1"][0] != "error message 1" {
		t.Fatalf("expected 'error message 1', got '%s'", e["field1"][0])
	}
	if e["field1"][1] != "error message 2" {
		t.Fatalf("expected 'error message 2', got '%s'", e["field1"][1])
	}
}

func TestErrorsGet(t *testing.T) {
	e := make(errors)

	got := e.Get("nonexistent")
	if got != "" {
		t.Fatalf("expected empty string, got '%s'", got)
	}

	e.Add("field1", "first error")
	e.Add("field1", "second error")

	got = e.Get("field1")
	if got != "first error" {
		t.Fatalf("expected 'first error', got '%s'", got)
	}
}

func TestErrorsGetReturnsOnlyFirst(t *testing.T) {
	e := make(errors)
	e.Add("field", "a")
	e.Add("field", "b")
	e.Add("field", "c")

	got := e.Get("field")
	if got != "a" {
		t.Fatalf("expected first error 'a', got '%s'", got)
	}
}
