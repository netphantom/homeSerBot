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
		queryRes := u.Db.Find(&notification, "user_id = ? AND process_id = ?", uid, pid)
		if queryRes.Error != nil {
			return nil, queryRes.Error
		}
		if (Notification{}) != notification {
			//Add it to the list that will be returned
			notificationList = append(notificationList, notification)
		}
	}
	return notificationList, nil
}

func (u *DbModel) RemoveNotification(n *Notification) error {
	queryRes := u.Db.Delete(&n)
	if queryRes.Error != nil {
		return queryRes.Error
	}
	return nil
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
		u.Db.Save(&n)
		return
	}
	u.Db.Create(&n)
}
