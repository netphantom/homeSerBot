package main

import (
	"flag"
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

func main() {
	addr := flag.String("addr", "4000", "The network port to use")
	ip := fmt.Sprintf("127.0.0.1:%s", *addr)

	dbUserName := flag.String("dbUserName", "admin", "The database username")
	dbPass := flag.String("dbPass", "admin", "The database password")
	dbIp := flag.String("dbIp", "8.8.8.8", "The database IP")
	dbPort := flag.Int("dbPort", 3306, "The database port")
	dbName := flag.String("dbName", "TelegramBot", "The database name")
	flag.Parse()

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	dsn := fmt.Sprint(*dbUserName, ":", *dbPass, "@tcp(", *dbIp, ":", *dbPort, ")/", *dbName, "?parseTime=true")
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
