package userauth

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
	"strings"
)

type UserInfo struct {
	UserId    string
	AvatarUrl string
}

func Signup(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	decoder := json.NewDecoder(r.Body)
	var newuser models.User
	if err := decoder.Decode(&newuser); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	var olduser models.User
	olduser.SelectByUserId(db, w, newuser.Userid)
	if olduser.Userid == newuser.Userid {
		httputils.RespondError(w, http.StatusInternalServerError, "userid is already used!")
	} else {
		newuser.Token = getToken(newuser.Userid)
		newuser.CreateUserRecord(db, w)
		httputils.RespondJson(w, http.StatusOK, newuser)
	}
}

func Signin(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	decoder := json.NewDecoder(r.Body)
	var newuser models.User
	if err := decoder.Decode(&newuser); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	var olduser models.User
	// olduser.SelectByUserId(w, newuser.Userid)
	if olduser.Userid != newuser.Userid {
		httputils.RespondError(w, http.StatusBadRequest, "user is not registered!")
	} else {
		if olduser.Password != newuser.Password {
			httputils.RespondError(w, http.StatusBadRequest, "password is different")
		} else {
			newuser.Token = olduser.Token
			httputils.RespondJson(w, http.StatusOK, newuser)
		}
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	token := r.Header.Get("Authorization")

	var userInfo UserInfo
	var user models.User
	var userFacePhoto models.UserFacePhoto

	user.GetUserInfoFromToken(db, w, token)

	db.Raw("SELECT * FROM user_face_photos WHERE userid = ?", user.Userid).Scan(&userFacePhoto)

	userInfo.UserId = user.Userid
	userInfo.AvatarUrl = userFacePhoto.Url

	httputils.RespondJson(w, http.StatusOK, userInfo)
}

func getToken(userid string) string {
	h := md5.New()
	h.Write([]byte(strings.ToLower(userid)))
	return hex.EncodeToString(h.Sum(nil))
}
