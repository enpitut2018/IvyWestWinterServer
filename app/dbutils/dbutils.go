package dbutils

import (
	"os"

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
	XID    string
	Userid string
	Url    string
}

type Download struct {
	gorm.Model
	Userid  string
	Photoid string
}

func ConnectPostgres() *gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL")) // osパッケージが必要
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
