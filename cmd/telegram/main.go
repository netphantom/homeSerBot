package telegram

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
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

func Telegram(tApi string, db gorm.DB) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	poller := &tb.LongPoller{Timeout: 10 * time.Second}

	b, err := tb.NewBot(tb.Settings{
		Token:  tApi,
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
		dbModel:    &mysqlmodels.DbModel{Db: &db},
	}

	bot.b.Handle("/start", bot.start)
	bot.b.Handle("/register", bot.register)

	bot.b.Handle("/pidList", bot.pidList)
	bot.b.Handle("/subscriptions", bot.subscriptions)
	bot.b.Handle("/subscribe", bot.subscribe)
	bot.b.Handle("/unsubscribe", bot.unsubscribe)

	go SendUpdates(&bot)

	bot.b.Start()
}
