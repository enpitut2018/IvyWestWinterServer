package userauth

import (
	"encoding/json"
	"net/http"

	"github.com/enpitut2018/IvyWestWinterServer/app/awsutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/enpitut2018/IvyWestWinterServer/app/faceidentification"
	"github.com/jinzhu/gorm"
	l "github.com/sirupsen/logrus"
)

type SourceRequest struct {
	Source string `json:"source"`
}

const s3FolderPath = "/user-face-photos"

func UploadUserFace(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	token := r.Header.Get("Authorization")
	var source SourceRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&source); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		l.Errorf(err.Error())
		return
	}

	var user models.User
	urlStr := awsutils.UploadPhoto(w, source.Source, s3FolderPath)
	if err := user.GetUserFromToken(db, token); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "No valid token.")
		l.Errorf("No valid token.")
		return
	}

	if err := user.UpdateAvatarURL(db, urlStr); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Cannot update AvatarURL.")
		l.Errorf("Cannot update AvatarURL.")
		return
	}
	httputils.RespondJson(w, http.StatusOK, user)
	faceidentification.AddUserFace(urlStr,user.AzurePersonID,w)
	l.Info("Success")
}
