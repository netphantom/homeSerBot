package mysqlmodels

import "gorm.io/gorm"

func (u *DbModel) UserProcessNotification(user *User) []Notification {
	var notificationList []Notification
	u.Db.Joins("JOIN processes ON notifications.process_id = processes.id AND sent = 0 AND user_id = ?", user.Id).Find(&notificationList)

	return notificationList
}

func (u *DbModel) MarkAsSent(n *Notification) {
	n.Sent = true
	u.Db.Save(&n)
}

func (u *DbModel) RemoveNotification(n *Notification) {
	u.Db.Delete(&n)
}

func (u *DbModel) AddNotification(n *Notification) {
	var lastNotification Notification
	queryRes := u.Db.Last(&lastNotification, "user_id = ? AND process_id = ?", n.UserID, n.ProcessID)
	if queryRes.Error != nil {
		if queryRes.Error == gorm.ErrRecordNotFound {
			u.Db.Create(&n)
			return
		}
	}
	if lastNotification.Active == n.Active && lastNotification.Process == n.Process {
		n.Model = lastNotification.Model
		n.ID = lastNotification.ID
		n.Sent = lastNotification.Sent
		u.Db.Save(&n)
		return
	}
	u.Db.Create(&n)
}
