package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserFacePhoto struct {
	gorm.Model
	XID    string
	Userid string
	Url    string
}