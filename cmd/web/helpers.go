package main

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"homeSerBot/pkg/mysqlmodels"
)

func (dash *dashboard) NotificationNumber(c *gin.Context) int {
	newUsers, err := dash.users.ListNewUsers()
	session := ginsession.FromContext(c)
	uid, ok := session.Get(sessionKey)
	if !ok {
		panic(err)
	}
	intUid := uid.(int)
	user := dash.users.VerifyId(uint(intUid))
	processNotification, err := dash.users.UserProcessNotification(user)

	if err != nil {
		panic(err)
	}
	return len(newUsers) + len(processNotification)
}

func ProcessInList(list []mysqlmodels.Process, value string) bool {
	for _, p := range list {
		if p.Name == value {
			return true
		}
	}
	return false
}
