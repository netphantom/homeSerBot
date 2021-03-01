package main

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"net/http"
)

func (dash *dashboard) AuthRequired(c *gin.Context) {
	session := ginsession.FromContext(c)
	_, ok := session.Get(sessionKey)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	c.Next()
}

func (dash *dashboard) UpdateNotificationNumber(c *gin.Context) {
	session := ginsession.FromContext(c)
	_, ok := session.Get("notifications")
	if !ok {
		session.Set("notifications", dash.NotificationNumber())
	}
	c.Next()
}
