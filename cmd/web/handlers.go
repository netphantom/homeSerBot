package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"homeSerBot/pkg/forms"
	"homeSerBot/pkg/mysqlmodels"
	"net/http"
)

const sessionKey = "uId"

func (dash *dashboard) showLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", &templateData{Form: forms.New(nil)})
}
func (dash *dashboard) showChangePassword(c *gin.Context) {
	c.HTML(http.StatusOK, "changePassword", gin.H{})
}

func (dash *dashboard) login(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		panic(err)
	}
	form := forms.New(c.Request.Form)
	form.Required("userName", "password")
	if !form.Valid() {
		form.Errors.Add("generic", "Please fill all of the fields")
		c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
		return
	}

	id, err := dash.users.Authenticate(form.Get("userName"), form.Get("password"))
	if err != nil {
		if errors.Is(err, mysqlmodels.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password are not correct")
			c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
		}
	} else {
		session := ginsession.FromContext(c)
		if form.Get("password") == "" {
			session.Set("init", true)
		} else {
			session.Set("init", false)
		}
		session.Set(sessionKey, id)
		err := session.Save()
		if err != nil {
			form.Errors.Add("generic", "Unable to create the sessions")
			c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
			return
		}
		c.Redirect(http.StatusFound, "/dashboard")
	}
}

func (dash *dashboard) logout(c *gin.Context) {
	session := ginsession.FromContext(c)
	_, ok := session.Get(sessionKey)
	form := forms.New(nil)
	if !ok {
		form.Errors.Add("generic", "Invalid session token")
		c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
		return
	}

	session.Delete(sessionKey)
	if err := session.Save(); err != nil {
		form.Errors.Add("generic", "Failed to delete the session")
		c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}

func (dash *dashboard) home(c *gin.Context) {
	session := ginsession.FromContext(c)
	init, ok := session.Get("init")
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	if init.(bool) {
		c.HTML(http.StatusOK, "changePassword", gin.H{
			"AlertMessage": "This is your first login. Please select your password",
		})
		return
	}
	c.HTML(http.StatusOK, "home", gin.H{})
}

func (dash *dashboard) profile(c *gin.Context) {
	session := ginsession.FromContext(c)
	uid, ok := session.Get(sessionKey)
	if !ok {
		c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
			"Error": "An error occurred during the session elaboration"})
		return
	}
	intUid := uid.(int)
	user := dash.users.VerifyId(uint(intUid))
	//TODO review the date
	c.HTML(http.StatusOK, "profile", gin.H{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	})
}

func (dash *dashboard) processes(c *gin.Context) {

}

func (dash *dashboard) changePassword(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		panic(err)
	}
	form := forms.New(c.Request.Form)
	form.Required("current", "new", "confirm")
	form.Minlength("new", 6)
	form.Minlength("confirm", 6)
	if form.Get("new") != form.Get("confirm") {
		form.Errors.Add("Alert", "Passwords do not match")
	}
	if !form.Valid() {
		c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
			"Alert": form.Errors.Get("Alert"),

			"Current": form.Errors.Get("current"),
			"New":     form.Errors.Get("new"),
			"Confirm": form.Errors.Get("Confirm"),
		})
		return
	}

	session := ginsession.FromContext(c)
	uid, ok := session.Get(sessionKey)
	if !ok {
		c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
			"Error": "An error occurred during the session elaboration"})
		return
	}

	init, ok := session.Get("init")
	if init.(bool) {
		err = dash.users.ChangePsw(form.Get("new"), form.Get(""), uid.(int))
		session.Delete("init")
	} else {
		err = dash.users.ChangePsw(form.Get("new"), form.Get("current"), uid.(int))
	}

	if err != nil {
		if errors.Is(err, mysqlmodels.ErrInvalidCredentials) {
			c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
				"Error": "Invalid Credentials"})
			return
		} else {
			panic(err)
		}
	}
	c.HTML(http.StatusOK, "home", gin.H{
		"Success": "Password Changed"})
}

func (dash *dashboard) showNotifications(c *gin.Context) {

	c.HTML(http.StatusOK, "notifications", gin.H{})
}

func (dash *dashboard) elaborateNotification(c *gin.Context) {

}
