package userauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/enpitut2018/IvyWestWinterServer/app/faceidentification"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	l "github.com/sirupsen/logrus"
)

type SignupRequest struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
}

func Signup(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	decoder := json.NewDecoder(r.Body)
	var requser SignupRequest
	if err := decoder.Decode(&requser); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		l.Errorf(err.Error())
		return
	}

	var user models.User
	user.SelectByUserID(db, requser.UserID)
	if user.UserID == requser.UserID {
		httputils.RespondError(w, http.StatusBadRequest, fmt.Sprintf("UserID %s is already exists.", requser.UserID))
		l.Errorf(fmt.Sprintf("UserID %s is already exists.", requser.UserID))
		return
	}

	user.UserID = requser.UserID
	user.Password = requser.Password
	// user.AzurePersonID = "0b4bbd63-ff70-423b-9aff-5263c745ff98" // 福山雅治の顔
	azurePersonID, err := faceidentification.CreatePerson(user.UserID, "My name is "+user.UserID, w)
	if err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		l.Errorf(err.Error())
		return
	}
	user.AzurePersonID = azurePersonID
	if err := user.CreateUserRecord(db); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		l.Errorf(err.Error())
		return
	}

	httputils.RespondJson(w, http.StatusOK, user)
	l.Infof("Success")
}

func Signin(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	decoder := json.NewDecoder(r.Body)
	var requser SignupRequest
	if err := decoder.Decode(&requser); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		l.Errorf(err.Error())
	}

	var user models.User
	user.SelectByUserID(db, requser.UserID)
	if user.UserID != requser.UserID {
		httputils.RespondError(w, http.StatusBadRequest, fmt.Sprintf("UserID [%s] is not found.", requser.UserID))
		l.Errorf(fmt.Sprintf("UserID [%s] is not found.", requser.UserID))
		return
	}

	if user.Password != requser.Password {
		httputils.RespondError(w, http.StatusBadRequest, "Password is different.")
		l.Errorf("Password is different.")
		return
	}

	httputils.RespondJson(w, http.StatusOK, user)
	l.Infof("Success.")
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user models.User
	token := r.Header.Get("Authorization")
	if err := user.GetUserFromToken(db, token); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Not valid token.")
		l.Errorf("Not valid token.")
		return
	}
	httputils.RespondJson(w, http.StatusOK, user)
	l.Infof("Success.")
}

type ResUser struct {
	UserID    string `gorm:"not null;unique"  json:"userid"`
	AvatarURL string `json:"avatarurl"`
}

func GetUsersInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var allUsers models.Users
	allUsers.GetAllUsers(db)
	queryString := r.URL.Query().Get("ids")

	var resUsers []ResUser
	if queryString == "" {
		for _, user := range allUsers.Users {
			var resUser ResUser
			resUser.UserID = user.UserID
			resUser.AvatarURL = user.AvatarURL
			resUsers = append(resUsers, resUser)
		}
	} else {
		for _, user := range allUsers.Users {
			for _, queryID := range strings.Split(queryString, ",") {
				if user.UserID == queryID {
					var resUser ResUser
					resUser.UserID = user.UserID
					resUser.AvatarURL = user.AvatarURL
					resUsers = append(resUsers, resUser)
				}
			}
		}
	}
	httputils.RespondJson(w, http.StatusOK, resUsers)
	return
}
