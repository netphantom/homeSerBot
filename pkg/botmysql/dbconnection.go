package botmysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDb(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	err = db.SetupJoinTable(&User{}, "Subscription",&UserProcess{})
	err = db.AutoMigrate(&User{}, &UserConnection{}, &Process{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
