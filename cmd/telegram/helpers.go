package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func IsAuthorized(bot *homeSerBot, u *tb.User) bool{
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
