package userauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	l "github.com/sirupsen/logrus"
)

type SignupRequest struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
}

func Signup(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		decoder := json.NewDecoder(c.Request().Body)
		var requser SignupRequest
		if err := decoder.Decode(&requser); err != nil {
			httputils.RespondError(w, http.StatusBadRequest, err.Error())
			l.Errorf(err.Error())
		}

		var user models.User
		user.SelectByUserID(db, requser.UserID)
		if user.UserID == requser.UserID {
			httputils.RespondError(w, http.StatusBadRequest, fmt.Sprintf("UserID %s is already exists.", requser.UserID))
			l.Errorf(fmt.Sprintf("UserID %s is already exists.", requser.UserID))
		} else {
			user.UserID = requser.UserID
			user.Password = requser.Password
			user.AzurePersonID = "0b4bbd63-ff70-423b-9aff-5263c745ff98" // 福山雅治の顔
			if ok := user.CreateUserRecord(db, w); ok {
				httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Success to create new user."})
				l.Infof("Success")
			}
		}
	}

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
	} else {
		if user.Password != requser.Password {
			httputils.RespondError(w, http.StatusBadRequest, "Password is different.")
			l.Errorf("Password is different.")
		} else {
			httputils.RespondJson(w, http.StatusOK, user)
			l.Infof("Success.")
		}
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user models.User
	token := r.Header.Get("Authorization")
	if ok := user.GetUserFromToken(db, w, token); ok {
		httputils.RespondJson(w, http.StatusOK, user)
		l.Infof("Success.")
	}
}
