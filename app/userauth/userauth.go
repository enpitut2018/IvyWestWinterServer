package userauth

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/dbutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"net/http"
	"strings"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newuser dbutils.User
	if err := decoder.Decode(&newuser); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	db := dbutils.ConnectPostgres()
	defer db.Close()

	// check already user
	var olduser dbutils.User
	db.Raw("SELECT * FROM USERS WHERE userid = ?", newuser.Userid).Scan(&olduser)
	if olduser.Userid == newuser.Userid {
		httputils.RespondError(w, http.StatusInternalServerError, "userid is already used!")
	} else {
		// create new user
		newuser.Token = getToken(newuser.Userid)
		if err := db.Create(&newuser).Error; err != nil {
			httputils.RespondError(w, http.StatusBadRequest, err.Error())
			panic(err)
		}
		httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "user created!"})
	}
}

func Signin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newuser dbutils.User
	if err := decoder.Decode(&newuser); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	db := dbutils.ConnectPostgres()
	defer db.Close()
	var olduser dbutils.User
	db.Raw("SELECT * FROM USERS WHERE userid = ?", newuser.Userid).Scan(&olduser)
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

func getToken(userid string) string {
	h := md5.New()
	h.Write([]byte(strings.ToLower(userid)))
	return hex.EncodeToString(h.Sum(nil))
}
