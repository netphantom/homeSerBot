package main

import (
	"flag"
	"fmt"
	"github.com/golangcollege/sessions"
	"homeSerBot/pkg/mysqlmodels"
	"log"
	"net/http"
	"os"
	"time"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type dashboard struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	users    *mysqlmodels.DbModel
}

func main() {
	addr := flag.String("addr", "4000", "The network port to use")
	secret := flag.String("secret", "secretValue", "Secret key")
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

	session := sessions.New([]byte(*secret))
	session.Lifetime = 24 * time.Hour

	dash := &dashboard{
		errorLog: errorLog,
		session:  session,
		users:    &mysqlmodels.DbModel{Db: db},
	}

	srv := &http.Server{
		Handler:      dash.routes(),
		Addr:         ip,
		ErrorLog:     errorLog,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
