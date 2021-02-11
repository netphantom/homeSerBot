package botmysql

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
	"time"
)

type DbModel struct {
	Db *gorm.DB
}

type User struct {
	gorm.Model
	tb.User
	Id			 uint `gorm:"primaryKey"`
	Allowed		 bool `gorm:"default:false"`
	UC     		 []UserConnection `gorm:"foreignKey:UID"`
	Subscription []Process `gorm:"many2many:user_process"`
}

type UserConnection struct {
	UID 		uint
	Command		string
	LastConnect time.Time
}

// Customization of the JoinTable for user_process
type UserProcess struct {
	gorm.Model
	UserID 		int `gorm:"primaryKey"`
	ProcessID 	int `gorm:"primaryKey"`
}

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
//Update the UserConnection with the current time and the command just sent
func (u *DbModel) UpdateLastInteraction(uid uint, command string) error {
	conn := UserConnection{
		UID:         uid,
		Command: 	 command,
		LastConnect: time.Now(),
	}
	queryResult := u.Db.Create(&conn)
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
	err:= u.Db.Model(&user).Association("Subscription").Find(&processList)
	if err != nil {
		return nil
	}
	return processList
}