package userauth

import (
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
)

type SignupRequest struct {
	Userid    string `json:"userid"`
	Password string `json:"password"`
}

func Signup(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	decoder := json.NewDecoder(r.Body)
	var requser SignupRequest
	if err := decoder.Decode(&requser); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
	}

	var user models.User
	user.SelectByUserId(db, requser.Userid)
	if user.Userid == requser.Userid {
		httputils.RespondError(w, http.StatusBadRequest, "Userid is already exists.")
	} else {
		user.Userid = requser.Userid
		user.Password = requser.Password
		if ok := user.CreateUserRecord(db, w); ok {
			httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Success to create new user."})
		}
	}
}

func Signin(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	decoder := json.NewDecoder(r.Body)
	var requser SignupRequest
	if err := decoder.Decode(&requser); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	var user models.User
	user.SelectByUserId(db, requser.Userid)
	if user.Userid != requser.Userid {
		httputils.RespondError(w, http.StatusBadRequest, "Userid is not found.")
	} else {
		if user.Password != requser.Password {
			httputils.RespondError(w, http.StatusBadRequest, "Password is different.")
		} else {
			httputils.RespondJson(w, http.StatusOK, user)
		}
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user models.User
	token := r.Header.Get("Authorization")
	if ok := user.GetUserFromToken(db, w, token);  !ok {
		httputils.RespondError(w, http.StatusUnauthorized, "not valid token.")
	} else {
		httputils.RespondJson(w, http.StatusOK, user)
	}
}
