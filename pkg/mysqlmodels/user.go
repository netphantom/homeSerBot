package mysqlmodels

import (
	"errors"
	"gorm.io/gorm"
)

//VerifyId function check if the user is correctly registered
func (u *DbModel) VerifyId(id uint) *User {
	var user User
	queryResult := u.Db.Find(&user, id)
	if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	if user.Allowed {
		return &user
	}
	return nil
}

func (u *DbModel) RegisterUser(user *User) error {
	queryResult := u.Db.Create(&user)
	if queryResult.Error != nil {
		return queryResult.Error
	}
	return nil
}

func (u *DbModel) SubscribeToProcess(user *User, pid string) (*Process, error) {
	process, err := u.GetProcessInfo(pid)
	if err != nil {
		return nil, err
	}
	err = u.Db.Model(&user).Association("Subscription").Append(process)
	if err != nil {
		return nil, err
	}
	return process, nil
}

func (u *DbModel) UnsubscribeToProcess(user *User, pid string) error {
	process, err := u.GetProcessInfo(pid)
	if err != nil {
		return err
	}
	err = u.Db.Model(&user).Association("Subscription").Delete(process)
	if err != nil {
		return err
	}
	return nil
}

func (u *DbModel) ListSubscribed(user *User) []Process {
	var processList []Process
	err := u.Db.Model(&user).Association("Subscription").Find(&processList)
	if err != nil {
		return nil
	}
	return processList
}
