package web

import (
	"github.com/gin-gonic/gin"
	"homeSerBot/pkg/mysqlmodels"
	"os"
	"os/exec"
	"strings"
	"time"
)

//NotificationNumber returns the number of notification to show
func (dash *dashboard) NotificationNumber(c *gin.Context) int {
	newUsers, err := dash.users.ListNewUsers()
	/*
		session := ginsession.FromContext(c)
		uid, ok := session.Get(sessionKey)
		if !ok {
			panic(err)
		}
		intUid := uid.(int)
		user := dash.users.VerifyId(uint(intUid))
		processNotification, err := dash.users.UserProcessNotification(user)
	*/
	if err != nil {
		panic(err)
	}
	return len(newUsers) //+ len(processNotification)
}

//UpdateProcessStatusUser updates the process status of the given user
func (dash *dashboard) UpdateProcessStatusUser(uid int) {
	processList, err := dash.users.ProcessList()
	if err != nil {
		panic(err)
	}

	for _, p := range processList {
		cmd := exec.Command("systemctl", "status", p.Name)
		cmd.Stderr = os.Stderr
		cmdOutput, _ := cmd.Output()

		fields := strings.Split(string(cmdOutput), "\n")

		var activeField string
		var statusField string
		for _, f := range fields {
			f = strings.TrimSpace(f)
			fieldSplit := strings.Split(f, ":")
			if fieldSplit[0] == "Active" {
				activeField = strings.Split(fieldSplit[1], "(")[0]
			} else if fieldSplit[0] == "Process" {
				statusField = strings.Split(fieldSplit[1], "status=")[1]
				statusField = strings.Trim(statusField, ")")
			}
		}

		notification := mysqlmodels.Notification{
			UserID:    uid,
			Process:   statusField,
			Active:    activeField,
			ProcessID: int(p.ID),
		}
		dash.users.AddNotification(&notification)
	}
}

//CreateNewNotifications is a process that every 5 minutes takes all the users registered in the dashboard and updates the status of their process
func CreateNewNotifications(dash *dashboard) {
	for {
		usersList, err := dash.users.ListAllUsers()
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Minute)
		for _, u := range usersList {
			uid := int(u.Id)
			go dash.UpdateProcessStatusUser(uid)
		}
	}
}

// ProcessInList checks if a given process name is present in a process list
func ProcessInList(list []mysqlmodels.Process, value string) bool {
	for _, p := range list {
		if p.Name == value {
			return true
		}
	}
	return false
}
