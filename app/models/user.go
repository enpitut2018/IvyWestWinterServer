package models

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strings"
)

type User struct {
	gorm.Model    `json:"-"`
	UserID        string `gorm:"not null;unique"  json:"userid"`
	AbaterURL     string `json:"abaterurl"`
	Password      string `json:"-"`
	Token         string `gorm:"not null;unique" json:"token"`
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

func (user *User) GetUserFromToken(db *gorm.DB, token string) error {
	if err := db.Find(&user, "token = ?", token).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) SelectByUserID(db *gorm.DB, userID string) error {
	if err := db.Find(&user, "user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) CreateUserRecord(db *gorm.DB) error {
	user.Token = getToken(user.UserID)
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) UpdateAbaterURL(db *gorm.DB, abaterurl string) error {
	user.AbaterURL = abaterurl
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (users *Users) GetAllUsers(db *gorm.DB) error {
	if err := db.Find(&users.Users).Error; err != nil {
		return err
	}
	return nil
}
