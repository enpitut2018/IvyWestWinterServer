package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Download struct {
	gorm.Model
	Userid   string
	PhotoUrl string
}