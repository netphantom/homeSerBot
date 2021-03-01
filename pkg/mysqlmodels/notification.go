package mysqlmodels

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
		if (Notification{}) == notification {
			return nil, ErrNoRecord
		}
		//Add it to the list that will be returned
		notificationList = append(notificationList, notification)
	}
	return notificationList, nil
}

func (u *DbModel) RemoveNotification(not *Notification) error {
	queryRes := u.Db.Delete(&not)
	if queryRes.Error != nil {
		return queryRes.Error
	}
	return nil
}
