package telegram

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

//IsAuthorized is a function that check if the bot has been authorized to communicate with the backend.
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

//EmptyPayload check if the content of a message is empty or not
func EmptyPayload(m *tb.Message) bool {
	if m.Payload == "" {
		return true
	}
	return false
}

//SendUpdates is a function that run on a different process and checks whether there are new updates to send to the user.
func SendUpdates(bot *homeSerBot) {
	for {
		time.Sleep(10 * time.Second)
		if bot.user != nil {
			notificationList := bot.dbModel.UserProcessNotification(bot.user)

			if len(notificationList) > 0 {
				for _, n := range notificationList {
					process, err := bot.dbModel.GetProcessInfo(n.ProcessID)
					if err != nil {
						panic(err)
					}
					message := fmt.Sprintf("Process: %s - Status: %s - Exit code: %s - On: %s", process.Name, n.Active, n.Process, n.UpdatedAt)
					bot.b.Send(bot.chat, message)
					bot.dbModel.MarkAsSent(&n)
				}
			}
		}
	}
}
