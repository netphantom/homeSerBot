package main

import (
	"flag"
	"fmt"
	"homeSerBot/cmd/telegram"
	"homeSerBot/cmd/web"
	"homeSerBot/pkg/mysqlmodels"
	"runtime"
)

func main() {
	tApi := flag.String("tApi", "<Telegram Token>", "The Telegram API token")

	addr := flag.String("addr", "4000", "The network port to use")
	ip := fmt.Sprintf("0.0.0.0:%s", *addr)

	dbUserName := flag.String("dbUserName", "admin", "The database username")
	dbPass := flag.String("dbPass", "admin", "The database password")
	dbIp := flag.String("dbIp", "8.8.8.8", "The database IP")
	dbPort := flag.Int("dbPort", 3306, "The database port")
	dbName := flag.String("dbName", "TelegramBot", "The database name")
	flag.Parse()

	dsn := fmt.Sprint(*dbUserName, ":", *dbPass, "@tcp(", *dbIp, ":", *dbPort, ")/", *dbName, "?parseTime=true")
	db, err := mysqlmodels.ConnectDb(dsn)
	if err != nil {
		panic(err)
	}

	go telegram.Telegram(*tApi, *db)
	go web.Web(ip, *db)

	runtime.Goexit()
}
