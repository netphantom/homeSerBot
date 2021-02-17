package main

import (
	"homeSerBot/pkg/forms"
	"homeSerBot/pkg/mysqlmodels"
	"time"
)

type templateData struct {
	Form *forms.Form
	User *mysqlmodels.User
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
