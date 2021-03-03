package telegram

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"homeSerBot/pkg/mysqlmodels"
	"strconv"
)

func (bot *homeSerBot) start(m *tb.Message) {
	message := "Welcome to the HomeServiceAlertBot.\nPlease register yourself using the /register command"
	bot.b.Send(m.Sender, message)
}

func (bot *homeSerBot) register(m *tb.Message) {
	userId := uint(m.Sender.ID)
	user := bot.dbModel.VerifyId(userId)
	if user == nil {
		newUser := mysqlmodels.User{
			User:         *m.Sender,
			Allowed:      false,
			Subscription: nil,
		}
		err := bot.dbModel.RegisterUser(&newUser)

		var message string
		if err == nil {
			message = "User not found, your details have been recorded for registration"
		} else {
			message = "Something weired happened during the registration process. The issue has been recorded"
			bot.infoLog.Print(err)
		}
		bot.b.Send(m.Sender, message)
	} else {
		if !user.Allowed {
			message := "Your request is still pending, please wait for authorization"
			bot.b.Send(m.Sender, message)
			return
		}
		bot.user = user
		bot.authorized = true
		bot.chat = m.Chat
		message := fmt.Sprintf("Welcome on board, %s", bot.user.Username)
		bot.b.Send(m.Sender, message)
	}
}

func (bot *homeSerBot) pidList(m *tb.Message) {
	if !IsAuthorized(bot, m.Sender) {
		return
	}

	pidList, err := bot.dbModel.ProcessList()
	if err != nil {
		bot.b.Send(m.Sender, err.Error())
	} else {
		message := "Here it comes the list with the process ID and the name"
		bot.b.Send(m.Sender, message)

		for _, p := range pidList {
			processInfo := fmt.Sprintf("Name: %s\nID: %d", p.Name, p.ID)
			bot.b.Send(m.Sender, processInfo)
		}
	}
}

func (bot *homeSerBot) subscribe(m *tb.Message) {
	if !IsAuthorized(bot, m.Sender) {
		return
	}

	if EmptyPayload(m) {
		message := "Please provide an argument"
		bot.b.Send(m.Sender, message)
		return
	}

	p := m.Payload
	pid, err := strconv.Atoi(p)
	if err != nil {
		bot.b.Send(m.Sender, "Please insert a number")
		return
	}

	process, err := bot.dbModel.SubscribeToProcess(bot.user, pid)
	if err != nil {
		bot.b.Send(m.Sender, err.Error())
		return
	}
	message := fmt.Sprintf("You have been subscribed to the process: %s, %d", process.Name, process.ID)
	bot.b.Send(m.Sender, message)
}

func (bot *homeSerBot) unsubscribe(m *tb.Message) {
	if !IsAuthorized(bot, m.Sender) {
		return
	}

	if EmptyPayload(m) {
		message := "Please provide an argument"
		bot.b.Send(m.Sender, message)
		return
	}
	p := m.Payload
	pid, err := strconv.Atoi(p)
	if err != nil {
		bot.b.Send(m.Sender, "Please insert a number")
		return
	}

	err = bot.dbModel.UnsubscribeToProcess(bot.user, pid)
	if err != nil {
		bot.b.Send(m.Sender, err.Error())
	} else {
		message := fmt.Sprintf("You have been unsubscribed from the process updates")
		bot.b.Send(m.Sender, message)
	}
}

func (bot *homeSerBot) subscriptions(m *tb.Message) {
	if !IsAuthorized(bot, m.Sender) {
		return
	}

	processList := bot.dbModel.ListSubscribed(bot.user)
	if processList == nil {
		bot.b.Send(m.Sender, "You do not have any subscriptions")
		return
	}
	bot.b.Send(m.Sender, "Please find below the process list you subscribed")
	for _, p := range processList {
		processInfo := fmt.Sprintf("Name: %s\nID: %d", p.Name, p.ID)
		bot.b.Send(m.Sender, processInfo)
	}
}

func (bot *homeSerBot) text(m *tb.Message) {
	bot.b.Send(m.Sender, m.Text)
}
