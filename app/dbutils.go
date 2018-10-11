package main

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Userid   string `gorm:"not null;unique"`
	Password string
	Token    string `gorm:"not null;unique"`
}

type Photo struct {
	gorm.Model
	Source   string
	Uploader User
}

type Download struct {
	gorm.Model
	User  User
	Photo Photo
}

func connectPostgres() *gorm.DB {
	db, err := gorm.Open("postgres",
		"host=db "+
			"port=5432 "+
			"user=postgres "+
			"sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func initialMigration() {
	db := connectPostgres()
	defer db.Close()
	db.AutoMigrate(&User{}, &Photo{}, &Download{})
}
