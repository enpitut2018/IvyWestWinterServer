package userauth

import (
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/awsutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	"net/http"
)

type SourceRequest struct {
	Source string `json:"source"`
}

func UploadUserFace(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	token := r.Header.Get("Authorization")

	var source SourceRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&source); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
	}

	s3FolderPath := "/user-face-photos"
	urlStr := awsutils.UploadPhoto(w, source.Source, s3FolderPath)

	var user models.User
	user.GetUserFromToken(db, w, token)
	user.UpdateAvatarURL(db, w, urlStr)
	httputils.RespondJson(w, http.StatusOK, user)
}
