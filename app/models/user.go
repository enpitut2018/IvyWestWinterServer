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
	Userid   string `gorm:"not null;unique"  json:"userid"`
	AbaterUrl string `json:"abaterurl"`
	Password string	`json:"-"`
	Token    string `gorm:"not null;unique" json:"token"`
}

func getToken(userid string) string {
	h := md5.New()
	h.Write([]byte(strings.ToLower(userid)))
	return hex.EncodeToString(h.Sum(nil))
}

func (user *User) GetUserFromToken(db *gorm.DB, w http.ResponseWriter, token string) bool {
	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, "not valid token.")
		return false
	} else {
		return true
	} 
}

func (user *User) SelectByUserId(db *gorm.DB, userid string) bool {
	if err := db.Raw("SELECT * FROM USERS WHERE userid = ?", userid).Scan(&user).Error; err != nil {
		// 見つからなかったらfalseを返す。
		return false
	} else {
		return true
	}
}

func (user *User) CreateUserRecord(db *gorm.DB, w http.ResponseWriter) bool {
	user.Token = getToken(user.Userid)
	if err := db.Create(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, err.Error())
		return false
	} else {
		return true
	}
}

func (user *User) UpdateAbaterUrl(db *gorm.DB, w http.ResponseWriter, abaterurl string) bool {
	user.AbaterUrl = abaterurl
	if err := db.Save(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, "Can't Update AbaterUrl")
		return false
	} else {
		return true
	}
}