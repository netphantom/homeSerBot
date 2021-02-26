package main

func (dash *dashboard) NotificationNumber() int {
	newUsers, err := dash.users.ListNewUsers()
	if err != nil {
		panic(err)
	}
	return len(newUsers)
}
