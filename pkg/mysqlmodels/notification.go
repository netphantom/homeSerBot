package mysqlmodels

import "gorm.io/gorm"

func (u *DbModel) UserProcessNotification(user *User) ([]Notification, error) {
	processList := u.ListSubscribed(user)
	uid := user.Id

	var notificationList []Notification
	for _, p := range processList {
		pid := p.ID

		//Find the notification in the table
		var notification Notification
		queryRes := u.Db.Find(&notification, "user_id = ? AND process_id = ? AND sent = 0", uid, pid)
		if queryRes.Error != nil {
			return nil, queryRes.Error
		}
		if notification.Active != "" {
			//Add it to the list that will be returned
			notificationList = append(notificationList, notification)
		}
	}
	return notificationList, nil
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
