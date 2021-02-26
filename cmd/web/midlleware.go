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
