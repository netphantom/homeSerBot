package main

import (
	"flag"
	"fmt"
	"homeSerBot/cmd/telegram"
	"homeSerBot/cmd/web"
	"runtime"
)

func main() {
	tApi := flag.String("tApi", "<Telegram Token>", "The Telegram API token")

	addr := flag.String("addr", "4000", "The network port to use")
	ip := fmt.Sprintf("127.0.0.1:%s", *addr)

	dbUserName := flag.String("dbUserName", "admin", "The database username")
	dbPass := flag.String("dbPass", "admin", "The database password")
	dbIp := flag.String("dbIp", "8.8.8.8", "The database IP")
	dbPort := flag.Int("dbPort", 3306, "The database port")
	dbName := flag.String("dbName", "TelegramBot", "The database name")
	flag.Parse()

	go telegram.Telegram(*dbUserName, *dbPass, *dbIp, *dbName, *tApi, *dbPort)
	go web.Web(ip, *dbUserName, *dbPass, *dbIp, *dbName, *dbPort)

	runtime.Goexit()
}
