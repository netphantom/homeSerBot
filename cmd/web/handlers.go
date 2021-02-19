package main

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"homeSerBot/pkg/forms"
	"letsgo/pkg/models"
	"net/http"
)

const sessionKey = "uId"

func (dash *dashboard) showLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.gohtml", &templateData{Form: forms.New(nil)})
}

func (dash *dashboard) login(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		panic(err)
	}
	form := forms.New(c.Request.Form)
	form.Required("userName", "password")
	if !form.Valid() {
		form.Errors.Add("generic", "Please provide the required data")
		c.HTML(http.StatusInternalServerError, "login.gohtml", &templateData{Form: form})
		return
	}

	session := sessions.Default(c)
	id, err := dash.users.Authenticate(form.Get("userName"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password are not correct")
			c.HTML(http.StatusInternalServerError, "login.gohtml", &templateData{Form: form})
		}
	} else {
		session.Set(sessionKey, id)
		if err := session.Save(); err != nil {
			form.Errors.Add("generic", "Unable to create the sessions")
			c.HTML(http.StatusInternalServerError, "login.gohtml", &templateData{Form: form})
			return
		}
		c.HTML(http.StatusOK, "index.html", nil)
	}
}
func (dash *dashboard) logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(sessionKey)
	form := forms.New(nil)
	if user == nil {
		form.Errors.Add("generic", "Invalid session token")
		c.HTML(http.StatusInternalServerError, "login.gohtml", &templateData{Form: form})
		return
	}

	session.Delete(sessionKey)
	if err := session.Save(); err != nil {
		form.Errors.Add("generic", "Failed to delete the session")
		c.HTML(http.StatusInternalServerError, "login.gohtml", &templateData{Form: form})
		return
	}
	c.HTML(http.StatusOK, "login.html", nil)
}
