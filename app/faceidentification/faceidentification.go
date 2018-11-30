package faceidentification

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	l "github.com/sirupsen/logrus"
)

var (
	personGroupID = "ivy-west-winter-test"
)

type FaceDetectRequest struct {
	URL string `json:"url"`
}

type FaceDetectResponse struct {
	FaceID         string
	FaceRectangle  interface{}
	FaceAttributes interface{}
}

type FaceVerifyRequest struct {
	FaceID        string `json:"faceId"`
	PersonID      string `json:"personId"`
	PersonGroupID string `json:"personGroupId"`
}

type FaceVerifyResponse struct {
	IsIdentical bool
	Confidence  float32
}

func PostAzureApi(url string, inJSON interface{}, outJSON interface{}, w http.ResponseWriter) bool {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(inJSON)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("AZURE_FACE_KEY"))
	client := &http.Client{}
	if res, _ := client.Do(req); res.StatusCode != 200 {
		httputils.RespondError(w, http.StatusBadRequest, res.Status)
		l.Errorf(res.Status)
		return false
	} else {
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&outJSON); err != nil {
			httputils.RespondError(w, http.StatusBadRequest, err.Error())
			l.Errorf(err.Error())
			return false
		} else {
			return true
		}
	}
}

func FaceDetect(url string, w http.ResponseWriter) (faceDetectResList []FaceDetectResponse, ok bool) {
	faceDetectURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/detect"
	inJSON := FaceDetectRequest{URL: url}
	PostAzureApi(faceDetectURL, inJSON, &faceDetectResList, w)
	return faceDetectResList, true
}

func FaceVerify(faceID string, personID string, w http.ResponseWriter) (faceVerifyRes FaceVerifyResponse, ok bool) {
	faceVerifyURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/verify"
	inJSON := FaceVerifyRequest{FaceID: faceID, PersonID: personID, PersonGroupID: personGroupID}
	PostAzureApi(faceVerifyURL, inJSON, &faceVerifyRes, w)
	return faceVerifyRes, true
}

func FaceIdentification(db *gorm.DB, w http.ResponseWriter, url string) ([]string, bool) {
	var allusers models.Users
	allusers.GetAllUsers(db, w)

	faceDetectResList, ok := FaceDetect(url, w)
	if !ok {
		return nil, false
	}

	l.Debugf("faceDetectResList: %+v\n\n", faceDetectResList)
	downloadUserIDs := make([]string, 0)
	for _, faceDetectRes := range faceDetectResList {
		for _, user := range allusers.Users {
			faceVerifyRes, ok := FaceVerify(faceDetectRes.FaceID, user.AzurePersonID, w)
			if !ok {
				return nil, false
			}
			l.Debugf("faceVerifyRes: %+v\n\n", faceVerifyRes)

			if faceVerifyRes.IsIdentical == true {
				download := models.Download{UserID: user.UserID, URL: url}
				if ok := download.CreateRecord(db, w); !ok {
					return nil, false
				}
				downloadUserIDs = append(downloadUserIDs, download.UserID)
			}
		}
	}
	return downloadUserIDs, true
}
