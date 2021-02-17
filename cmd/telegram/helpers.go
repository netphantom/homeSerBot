package telegram

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"homeSerBot/pkg/mysqlmodels"
	"time"
)

func IsAuthorized(bot *homeSerBot, u *tb.User) bool {
	if !bot.authorized {
		message := "You are not authorized to perform this action"
		bot.b.Send(u, message)

		if u.IsBot {
			bot.infoLog.Printf("A bot named %s tried to to use private functions", u.Username)
		} else {
			bot.infoLog.Printf("A human named %s tried to to use private functions", u.Username)
		}

		return false
	}
	return true
}

func EmptyPayload(m *tb.Message) bool {
	if m.Payload == "" {
		return true
	}
	return false
}

func SendUpdates(bot *homeSerBot) {
	for {
		time.Sleep(3 * time.Second)
		if bot.user != nil {
			notificationList, err := bot.dbModel.GetUserNotification(bot.user)
			if err != nil {
				if errors.Is(err, mysqlmodels.ErrNoRecord) {
					continue
				}
				panic(err)
			}
			if len(notificationList) > 0 {
				for _, n := range notificationList {
					bot.b.Send(bot.chat, n.Status)
					err = bot.dbModel.RemoveNotification(&n)
					if err != nil {
						bot.b.Send(bot.chat, err)
					}
				}
			}
		}
	}
}
