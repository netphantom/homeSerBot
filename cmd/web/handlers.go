package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	tb "gopkg.in/tucnak/telebot.v2"
	"homeSerBot/pkg/forms"
	"homeSerBot/pkg/mysqlmodels"
	"net/http"
)

const sessionKey = "uId"

func (dash *dashboard) showLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", &templateData{Form: forms.New(nil)})
}
func (dash *dashboard) showChangePassword(c *gin.Context) {
	sess := ginsession.FromContext(c)
	notifications, _ := sess.Get("notifications")

	c.HTML(http.StatusOK, "changePassword", gin.H{
		"Notifications": notifications,
	})
}
func (dash *dashboard) showNotifications(c *gin.Context) {
	newUsersList, err := dash.users.ListNewUsers()
	if err != nil {
		panic(err)
	}
	var newUsers []tb.User
	for _, u := range newUsersList {
		newUsers = append(newUsers, u.User)
	}
	sess := ginsession.FromContext(c)
	sess.Set("notifications", 0)
	c.HTML(http.StatusOK, "notifications", gin.H{
		"NewUsers": newUsers,
	})
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
		sess := ginsession.FromContext(c)
		if form.Get("password") == "" {
			sess.Set("init", true)
		} else {
			sess.Set("init", false)
		}
		sess.Set(sessionKey, id)
		err := sess.Save()
		if err != nil {
			form.Errors.Add("generic", "Unable to create the sessions")
			c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
			return
		}
		c.Redirect(http.StatusFound, "/dashboard")
	}
}

func (dash *dashboard) logout(c *gin.Context) {
	sess := ginsession.FromContext(c)
	_, ok := sess.Get(sessionKey)
	form := forms.New(nil)
	if !ok {
		form.Errors.Add("generic", "Invalid sess token")
		c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
		return
	}

	sess.Delete(sessionKey)
	if err := sess.Save(); err != nil {
		form.Errors.Add("generic", "Failed to delete the sess")
		c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}

func (dash *dashboard) home(c *gin.Context) {
	sess := ginsession.FromContext(c)
	init, ok := sess.Get("init")
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
	notifications, _ := sess.Get("notifications")
	c.HTML(http.StatusOK, "home", gin.H{
		"Notifications": notifications,
	})
}

func (dash *dashboard) profile(c *gin.Context) {
	sess := ginsession.FromContext(c)
	uid, ok := sess.Get(sessionKey)
	notifications, _ := sess.Get("notifications")
	if !ok {
		c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
			"Notifications": notifications,
			"Error":         "An error occurred during the sess elaboration"})
		return
	}
	intUid := uid.(int)
	user := dash.users.VerifyId(uint(intUid))
	// TODO review the date
	c.HTML(http.StatusOK, "profile", gin.H{
		"Notifications": notifications,

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
	sess := ginsession.FromContext(c)
	notifications, _ := sess.Get("notifications")
	if !form.Valid() {
		c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
			"Alert":         form.Errors.Get("Alert"),
			"Notifications": notifications,

			"Current": form.Errors.Get("current"),
			"New":     form.Errors.Get("new"),
			"Confirm": form.Errors.Get("Confirm"),
		})
		return
	}

	uid, ok := sess.Get(sessionKey)
	if !ok {
		c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
			"Notifications": notifications,

			"Error": "An error occurred during the sess elaboration"})
		return
	}

	init, ok := sess.Get("init")
	if init.(bool) {
		err = dash.users.ChangePsw(form.Get("new"), form.Get(""), uid.(int))
		sess.Delete("init")
	} else {
		err = dash.users.ChangePsw(form.Get("new"), form.Get("current"), uid.(int))
	}

	if err != nil {
		if errors.Is(err, mysqlmodels.ErrInvalidCredentials) {
			c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
				"Notifications": notifications,
				"Error":         "Invalid Credentials"})
			return
		} else {
			panic(err)
		}
	}
	c.HTML(http.StatusOK, "home", gin.H{
		"Notifications": notifications,

		"Success": "Password Changed"})
}

func (dash *dashboard) adminMode(c *gin.Context) {
	username := c.Query("username")
	status := c.Query("status")
	if status == "accept" {
		err := dash.users.AllowUser(username)
		if err != nil {
			c.HTML(http.StatusOK, "notifications", gin.H{
				"Error": "An error occurred during the acceptance process",
			})
			return
		}
	} else {
		err := dash.users.NotAllowUser(username)
		if err != nil {
			c.HTML(http.StatusOK, "notifications", gin.H{
				"Error": "An error occurred during the acceptance process",
			})
			return
		}
	}
	c.HTML(http.StatusOK, "notifications", gin.H{
		"Success": "Operation performed correctly",
	})
}
