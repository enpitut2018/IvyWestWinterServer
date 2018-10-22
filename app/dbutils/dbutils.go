package dbutils

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Userid   string `gorm:"not null;unique"`
	Password string
	Token    string `gorm:"not null;unique"`
}

type Photo struct {
	gorm.Model
	Source string
	Userid string
	Url    string
}

type Download struct {
	gorm.Model
	User  User
	Photo Photo
}

func ConnectPostgres() *gorm.DB {
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

func InitialMigration() {
	db := ConnectPostgres()
	defer db.Close()
	db.AutoMigrate(&User{}, &Photo{}, &Download{})
}

func InitialSQL() {
	db := ConnectPostgres()
	defer db.Close()

	user := User{Userid: "ivy", Password: "pass", Token: "AAAAAAAA"}
	if err := db.Create(&user).Error; err != nil {
		panic(err.Error())
	}

	photo := Photo{Source: "ENCODEDPHOTO", Userid: user.Userid}
	if err := db.Create(&photo).Error; err != nil {
		panic(err.Error())
	}
}
