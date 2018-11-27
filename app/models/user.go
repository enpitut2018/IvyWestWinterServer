package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"net/http"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
)

type User struct {
	gorm.Model `json:"-"`
	UserID   string `gorm:"not null;unique"  json:"userid"`
	AbaterURL string `json:"abaterurl"`
	Password string	`json:"-"`
	Token    string `gorm:"not null;unique" json:"token"`
	AzurePersonID string `json:"-"`
}

type Users struct {
	Users []User
}

func getToken(userID string) string {
	h := md5.New()
	h.Write([]byte(strings.ToLower(userID)))
	return hex.EncodeToString(h.Sum(nil))
}

func (user *User) GetUserFromToken(db *gorm.DB, w http.ResponseWriter, token string) bool {
	if err := db.Find(&user, "token = ?", token).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, "Not valid token.")
		panic("Not valid token.")
	}
	return true 
}

func (user *User) SelectByUserID(db *gorm.DB, userID string) bool {
	if err := db.Find(&user, "user_id = ?", userID).Error; err != nil {
		return false
	}
	return true
}

func (user *User) CreateUserRecord(db *gorm.DB, w http.ResponseWriter) bool {
	user.Token = getToken(user.UserID)
	if err := db.Create(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, err.Error())
		panic(err.Error())
		return false
	} else {
		return true
	}
}

func (user *User) UpdateAbaterURL(db *gorm.DB, w http.ResponseWriter, abaterurl string) bool {
	user.AbaterURL = abaterurl
	if err := db.Save(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, "Can't Update AbaterURL.")
		panic("Can't Update AbaterURL.")
		return false
	}
	return true

}

func (users *Users) GetAllUsers(db *gorm.DB, w http.ResponseWriter) bool {
	if err := db.Find(&users.Users).Error; err != nil{
		return false
	}
	return true
}