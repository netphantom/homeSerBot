package telegram

import (
	"flag"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"homeSerBot/pkg/mysqlmodels"
	"log"
	"os"
	"time"
)

type homeSerBot struct {
	b          *tb.Bot
	dbModel    *mysqlmodels.DbModel
	authorized bool
	infoLog    *log.Logger
	user       *mysqlmodels.User
	chat       *tb.Chat
}

func main() {
	tApi := flag.String("tApi", "<Telegram Token>", "The Telegram API token")
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

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	poller := &tb.LongPoller{Timeout: 10 * time.Second}

	b, err := tb.NewBot(tb.Settings{
		Token:  *tApi,
		Poller: poller,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	bot := homeSerBot{
		b:          b,
		authorized: false,
		infoLog:    infoLog,
		dbModel:    &mysqlmodels.DbModel{Db: db},
	}

	bot.b.Handle("/start", bot.start)
	bot.b.Handle("/register", bot.register)

	bot.b.Handle("/pidList", bot.pidList)
	bot.b.Handle("/subscriptions", bot.subscriptions)
	bot.b.Handle("/subscribe", bot.subscribe)
	bot.b.Handle("/unsubscribe", bot.unsubscribe)
	bot.b.Handle(tb.OnText, bot.text)

	go SendUpdates(&bot)

	bot.b.Start()
}
