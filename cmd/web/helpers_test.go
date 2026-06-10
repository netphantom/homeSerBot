package web

import (
	"testing"

	"homeSerBot/pkg/mysqlmodels"
)

func TestProcessInList_Found(t *testing.T) {
	list := []mysqlmodels.Process{
		{Name: "nginx.service"},
		{Name: "postgresql.service"},
	}

	if !ProcessInList(list, "nginx.service") {
		t.Fatal("expected to find nginx.service in list")
	}
}

func TestProcessInList_NotFound(t *testing.T) {
	list := []mysqlmodels.Process{
		{Name: "nginx.service"},
	}

	if ProcessInList(list, "apache.service") {
		t.Fatal("expected not to find apache.service in list")
	}
}

func TestProcessInList_EmptyList(t *testing.T) {
	list := []mysqlmodels.Process{}

	if ProcessInList(list, "anything.service") {
		t.Fatal("expected not to find anything in empty list")
	}
}

func TestProcessInList_MultipleMatches(t *testing.T) {
	list := []mysqlmodels.Process{
		{Name: "a.service"},
		{Name: "b.service"},
		{Name: "a.service"},
	}

	if !ProcessInList(list, "a.service") {
		t.Fatal("expected to find a.service in list with duplicates")
	}
}
