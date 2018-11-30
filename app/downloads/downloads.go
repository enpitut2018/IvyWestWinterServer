package downloads

import (
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	l "github.com/sirupsen/logrus"
	"net/http"
)

func CreateDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}

func GetDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user models.User
	token := r.Header.Get("Authorization")
	if err := user.GetUserFromToken(db, token); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Not valid token.")
		l.Errorf("Not valid token.")
		return
	}

	var downloads models.Downloads
	if err := downloads.GetDownloadsByUserID(db, user.UserID); err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, err.Error())
		l.Errorf(err.Error())
		return
	}

	httputils.RespondJson(w, http.StatusOK, downloads.Downloads)
	l.Infof("Success.")
}

func DeleteDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}
