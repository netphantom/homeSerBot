package web

import (
	"errors"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	tb "gopkg.in/tucnak/telebot.v2"
	"homeSerBot/pkg/forms"
	"homeSerBot/pkg/mysqlmodels"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

const sessionKey = "uId"

//showLogin renders the login page to the user
func (dash *dashboard) showLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", &templateData{Form: forms.New(nil)})
}

//showChangePassword renders the page that allow the user to change the password
func (dash *dashboard) showChangePassword(c *gin.Context) {
	sess := ginsession.FromContext(c)
	notifications, _ := sess.Get("notifications")

	c.HTML(http.StatusOK, "changePassword", gin.H{
		"Notifications": notifications,
	})
}

//showNotifications is a function that renders the notification number on the top of the navigation bar
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

	uid, ok := sess.Get(sessionKey)
	if !ok {
		panic(err)
	}
	intUid := uid.(int)
	user := dash.users.VerifyId(uint(intUid))
	processNotifList := dash.users.UserProcessNotification(user)

	var processesNotification []mysqlmodels.Notification
	for _, p := range processNotifList {
		processesNotification = append(processesNotification, p)
	}
	c.HTML(http.StatusOK, "notifications", gin.H{
		"NewUsers":    newUsers,
		"UserProcess": processesNotification,
	})

}

//showProcessesList shows the page with the list of the processes available for the registration
func (dash *dashboard) showProcessesList(c *gin.Context) {
	processList, err := dash.users.ProcessList()
	sess := ginsession.FromContext(c)
	notifications, _ := sess.Get("notifications")

	if err != nil {
		c.HTML(http.StatusOK, "process", gin.H{
			"Notifications": notifications,

			"Error": "Something went wrong during the process list elaboration",
		})
		return
	}

	c.HTML(http.StatusOK, "process", gin.H{
		"Notifications": notifications,

		"Process": processList,
	})
}

//showProcessDetails shows the details of the selected process
func (dash *dashboard) showProcessDetails(c *gin.Context) {
	p := c.Query("process")
	sess := ginsession.FromContext(c)
	notifications, _ := sess.Get("notifications")

	pid, err := strconv.Atoi(p)
	if err != nil {
		c.HTML(http.StatusOK, "process", gin.H{
			"Notifications": notifications,

			"Error": "Something went wrong",
		})
		return
	}

	process, err := dash.users.GetProcessInfo(pid)
	if err != nil {
		c.HTML(http.StatusOK, "process", gin.H{
			"Notifications": notifications,

			"Error": "Something went wrong",
		})
		return
	}
	c.HTML(http.StatusOK, "processDetail", gin.H{
		"Notifications": notifications,

		"Process": process,
	})
}

//showNewProcess renders the page showing all the process available for the registration by performing an ls command and parsing the output
func (dash *dashboard) showNewProcess(c *gin.Context) {
	cmdOutput, err := exec.Command("ls", "/etc/systemd/system/").Output()
	sess := ginsession.FromContext(c)
	notifications, _ := sess.Get("notifications")

	if err != nil {
		c.HTML(http.StatusOK, "processAdd", gin.H{
			"Notifications": notifications,

			"Error": err,
		})
		return
	}
	fileList := strings.Split(string(cmdOutput), "\n")
	dbProcessList, err := dash.users.ProcessList()
	if err != nil {
		c.HTML(http.StatusOK, "processAdd", gin.H{
			"Notifications": notifications,

			"Error": err,
		})
		return
	}

	var purgedFileList []string
	for _, f := range fileList {
		if strings.HasSuffix(f, ".service") {
			if !ProcessInList(dbProcessList, f) {
				purgedFileList = append(purgedFileList, f)
			}
		}
	}

	c.HTML(http.StatusOK, "processAdd", gin.H{
		"Notifications": notifications,

		"Process": purgedFileList,
	})
}

//login performs the user login into the dashboard
func (dash *dashboard) login(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		panic(err)
	}

	form := forms.New(c.Request.Form)
	id, err := dash.users.Authenticate(form.Get("userName"), form.Get("password"))
	user := dash.users.VerifyId(uint(id))
	if err != nil {
		//The user has been found in the DB
		form.Required("userName", "password")
		if !form.Valid() {
			form.Errors.Add("generic", "Please fill all of the fields")
			c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
			return
		}

		if errors.Is(err, mysqlmodels.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password are not correct")
			c.HTML(http.StatusInternalServerError, "login", &templateData{Form: form})
		}
	} else {
		//The user is new, thus it's needed to register him/her
		sess := ginsession.FromContext(c)
		if user.Password == nil {
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

//logout performs the user logout
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

//home renders the home page, showing some process upd ates
func (dash *dashboard) home(c *gin.Context) {
	sess := ginsession.FromContext(c)
	init, ok := sess.Get("init")
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	if init.(bool) {
		c.HTML(http.StatusOK, "changePassword", gin.H{
			"Alert": "This is your first login. Please select your password",
		})
		return
	}
	notifications, _ := sess.Get("notifications")
	uid, ok := sess.Get(sessionKey)
	if !ok {
		c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
			"Notifications": notifications,
			"Error":         "An error occurred during the session elaboration"})
		return
	}

	intUid := uid.(int)
	user := dash.users.VerifyId(uint(intUid))
	userNotification := dash.users.UserProcessNotification(user)

	c.HTML(http.StatusOK, "home", gin.H{
		"userNot":       userNotification,
		"Notifications": notifications,
	})
}

//profile shows the user profile with the details
func (dash *dashboard) profile(c *gin.Context) {
	sess := ginsession.FromContext(c)
	uid, ok := sess.Get(sessionKey)
	notifications, _ := sess.Get("notifications")
	if !ok {
		c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
			"Notifications": notifications,
			"Error":         "An error occurred during the session elaboration"})
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

//changePassword allows the user to change the password
func (dash *dashboard) changePassword(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		panic(err)
	}
	sess := ginsession.FromContext(c)
	init, ok := sess.Get("init")
	form := forms.New(c.Request.Form)

	if !init.(bool) {
		form.Required("current", "new", "confirm")
	} else {
		form.Required("new", "confirm")
	}
	form.Minlength("new", 6)
	form.Minlength("confirm", 6)
	if form.Get("new") != form.Get("confirm") {
		form.Errors.Add("Alert", "Passwords do not match")
	}
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

			"Error": "An error occurred during the session elaboration"})
		return
	}

	if init.(bool) {
		err = dash.users.ChangePsw(form.Get("new"), "", uid.(int))
		user := dash.users.VerifyId(uint(uid.(int)))
		_ = dash.users.AllowUser(user.Username)
		notifications = 0
		sess.Set("init", false)
	} else {
		err = dash.users.ChangePsw(form.Get("new"), form.Get("current"), uid.(int))
	}

	if err != nil {
		if errors.Is(err, mysqlmodels.ErrInvalidCredentials) {
			c.HTML(http.StatusInternalServerError, "changePassword", gin.H{
				"Notifications": notifications,

				"Error": "Invalid Credentials"})
			return
		} else {
			panic(err)
		}
	}
	c.HTML(http.StatusOK, "home", gin.H{
		"Notifications": notifications,

		"Success": "Password Changed"})
}

//adminMode allows the user to accept or not new users
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

//deleteProcess remove a process from the "available" list of the user
func (dash *dashboard) deleteProcess(c *gin.Context) {
	p := c.Query("process")
	sess := ginsession.FromContext(c)
	notifications, _ := sess.Get("notifications")

	pid, err := strconv.Atoi(p)
	if err != nil {
		c.HTML(http.StatusOK, "process", gin.H{
			"Notifications": notifications,

			"Error": "Something went wrong",
		})
		return
	}
	err = dash.users.DeleteProcess(pid)
	if err != nil {
		c.HTML(http.StatusOK, "process", gin.H{
			"Notifications": notifications,

			"Error": "Something went wrong",
		})
		return
	}
	processList, err := dash.users.ProcessList()
	if err != nil {
		c.HTML(http.StatusOK, "process", gin.H{
			"Notifications": notifications,
			"Process":       processList,

			"Error": "Something went wrong during the process list elaboration",
		})
		return
	}

	c.HTML(http.StatusOK, "process", gin.H{
		"Notifications": notifications,
		"Process":       processList,

		"Success": "Operation performed correctly",
	})

}

//editProcess allows the user to edit the process description
func (dash *dashboard) editProcess(c *gin.Context) {
	err := c.Request.ParseForm()
	sess := ginsession.FromContext(c)
	notifications, _ := sess.Get("notifications")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "processDetail", gin.H{
			"Error":         err,
			"Notifications": notifications,
		})
		return
	}

	p := c.Request.Form.Get("pid")
	description := c.Request.Form.Get("description")
	pid, err := strconv.Atoi(p)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "processDetail", gin.H{
			"Error":         err,
			"Notifications": notifications,
		})
		return
	}

	err = dash.users.UpdateDescription(pid, description)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "processDetail", gin.H{
			"Error":         err,
			"Notifications": notifications,
		})
		return
	}
	processList, err := dash.users.ProcessList()
	c.HTML(http.StatusOK, "process", gin.H{
		"Success":       "Description updated",
		"Process":       processList,
		"Notifications": notifications,
	})
}

//processAdd adds a process to the "available" list of the user
func (dash *dashboard) processAdd(c *gin.Context) {
	err := c.Request.ParseForm()
	sess := ginsession.FromContext(c)
	notifications, _ := sess.Get("notifications")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "processAdd", gin.H{
			"Error": err,

			"Notifications": notifications,
		})
		return
	}
	description := c.Request.Form["description"]
	pidName := c.Request.Form["pidName"]

	for i, d := range description {
		if d == "" {
			continue
		}
		dash.users.AddProcess(pidName[i], d)
	}

	processList, err := dash.users.ProcessList()
	c.HTML(http.StatusOK, "process", gin.H{
		"Success": "Processes successfully added",
		"Process": processList,

		"Notifications": notifications,
	})
}

//notFound handles the event of a requested page not found
func (dash *dashboard) notFound(c *gin.Context) {
	sess := ginsession.FromContext(c)
	uid, _ := sess.Get(sessionKey)
	if uid != nil {
		notifications, _ := sess.Get("notifications")
		c.HTML(http.StatusNotFound, "home", gin.H{
			"Notifications": notifications,

			"Error": "The page requested cannot be found ",
		})
	} else {
		c.HTML(http.StatusNotFound, "notFound", nil)
	}
}
