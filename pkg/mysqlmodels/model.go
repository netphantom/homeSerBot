package mysqlmodels

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var (
	ErrNoRecord           = errors.New("model: No matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
)

type DbModel struct {
	Db *gorm.DB
}

type Process struct {
	gorm.Model
	Name        string
	Description string
}

type Notification struct {
	gorm.Model
	UserID    int `gorm:"primaryKey"`
	ProcessID int `gorm:"primaryKey"`
	Active    string
	Process   string
}

type User struct {
	gorm.Model
	tb.User
	Id           uint      `gorm:"primaryKey"`
	Allowed      bool      `gorm:"default:false"`
	Password     []byte    `gorm:"default:Null"`
	Subscription []Process `gorm:"many2many:user_process"`
	Notification []Notification
}

// Customization of the JoinTable for user_process
type UserProcess struct {
	gorm.Model
	UserID    int `gorm:"primaryKey"`
	ProcessID int `gorm:"primaryKey"`
}
