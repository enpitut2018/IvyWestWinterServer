package user

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/enpitut2018/IvyWestWinterServer/app/dbutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"net/http"
)

type User struct {
	gorm.Model
	Userid   string `gorm:"not null;unique"`
	AbaterUrl string
	Password string
	Token    string `gorm:"not null;unique"`
}

func (u *User) GetUserInfoFromToken(db *gorm.DB, token string) {
	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, err.Error())
		panic(err.Error())
	}
}