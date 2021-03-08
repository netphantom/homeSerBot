package web

import (
	"fmt"
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

func Web(ip, dbUserName, dbPass, dbIp, dbName string, dbPort int) {

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	dsn := fmt.Sprint(dbUserName, ":", dbPass, "@tcp(", dbIp, ":", dbPort, ")/", dbName, "?parseTime=true")
	db, err := mysqlmodels.ConnectDb(dsn)
	if err != nil {
		panic(err)
	}

	dash := &dashboard{
		errorLog: errorLog,
		users:    &mysqlmodels.DbModel{Db: db},
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
