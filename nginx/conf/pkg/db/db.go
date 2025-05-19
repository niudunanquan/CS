package db

import (
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var sqlDrivers = map[string]func(string) gorm.Dialector{
	"postgres": postgres.Open,
	"mysql":    mysql.Open,
	"sqlite":   sqlite.Open,
}

type Model struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at"`
}

var Db *gorm.DB

func InitDB(driver, dsn string) error {
	var err error

	Db, err = gorm.Open(sqlDrivers[driver](dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		return err
	}

	db, err := Db.DB()
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	err = Db.AutoMigrate(&User{}, &Site{}, &Cert{}, &Upstream{})
	if err != nil {
		return err
	}

	var user User
	count, err := user.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		user.Username = "admin"
		user.Password = "Passw0rd!"
		user.Email = "support@uusec.com"
		err = user.Insert()
	}

	return err
}

func Close() error {
	db, err := Db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
