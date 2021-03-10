package web

import (
	"gorm.io/gorm"
	"homeSerBot/pkg/mysqlmodels"
	"log"
	"net/http"
	"os"
	"time"
)

type dashboard struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    *mysqlmodels.DbModel
}

func Web(ip string, db gorm.DB) {
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	dash := &dashboard{
		errorLog: errorLog,
		users:    &mysqlmodels.DbModel{Db: &db},
	}

	go CreateNewNotifications(dash)

	srv := &http.Server{
		Handler:      dash.routes(),
		Addr:         ip,
		ErrorLog:     errorLog,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
