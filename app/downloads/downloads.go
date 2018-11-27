package downloads

import (
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"net/http"
	"github.com/jinzhu/gorm"
)

func CreateDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}

func GetDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user models.User
	token := r.Header.Get("Authorization")
	user.GetUserFromToken(db, w, token)
	
	var downloads models.Downloads
	downloads.GetDownloadsByUserID(db, w, user.UserID)
	httputils.RespondJson(w, http.StatusOK, downloads.Downloads)
}

func DeleteDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}
