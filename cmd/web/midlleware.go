package main

import (
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"net/http"
)

//AuthRequired checks if the user id has been cached in the gin context
func (dash *dashboard) AuthRequired(c *gin.Context) {
	session := ginsession.FromContext(c)
	_, ok := session.Get(sessionKey)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	c.Next()
}

//SetNotificationNumber set the entry of the notification number in the gin context
func (dash *dashboard) SetNotificationNumber(c *gin.Context) {
	session := ginsession.FromContext(c)
	_, ok := session.Get("notifications")
	if !ok {
		session.Set("notifications", dash.NotificationNumber(c))
	}
	c.Next()
}
