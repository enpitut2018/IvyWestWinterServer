package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"net/http"
)

type User struct {
	gorm.Model
	Userid   string `gorm:"not null;unique"  json:"id"`
	AbaterUrl string `json:"id"`
	Password string	`json:"-"`
	Token    string `gorm:"not null;unique"`
}

func (user *User) GetUserInfoFromToken(db *gorm.DB, w http.ResponseWriter, token string) {
	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, "Can't get User From Token")
		panic(err.Error())
	}
}

func (user *User) SelectByUserId(db *gorm.DB, w http.ResponseWriter, userid string) {
	if err := db.Raw("SELECT * FROM USERS WHERE userid = ?", userid).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, "already userid")
		panic("already userid\n")
	}
}

func (user *User) CreateUserRecord(db *gorm.DB, w http.ResponseWriter) {
	if err := db.Create(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "User Create Record Error")
		panic("User Create Record Error\n")
	}
}
