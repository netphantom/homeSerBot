package mysqlmodels

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"letsgo/pkg/models"
)

//VerifyId function check if the user is correctly registered
func (u *DbModel) VerifyId(id uint) *User {
	var user User
	queryResult := u.Db.Find(&user, id)
	if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return &user
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

func (u *DbModel) Authenticate(username, password string) (int, error) {
	var user User
	queryResult := u.Db.First(&user, "username = ?", username)
	if errors.Is(queryResult.Error, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	if string(user.Password) == "" {
		return int(user.Id), nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return int(user.Id), nil
}
