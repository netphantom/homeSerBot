package mysqlmodels

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"

	"gorm.io/gorm"
)

func ConnectDb(dsn string, dbType string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch dbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	default:
		err = errors.New("database dialect not recognized")
	}

	if err != nil {
		return nil, err
	}

	err = db.SetupJoinTable(&User{}, "Subscription", &UserProcess{})
	err = db.AutoMigrate(&User{}, &Process{}, &Notification{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
