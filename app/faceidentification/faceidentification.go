package faceidentification

import (
	"bytes"
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"

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

type CreatePersonRequest struct {
	Name     string `json:name`
	UserData string `json:userData`
}

type CreatePersonResponse struct {
	PersonID string
}

func PostAzureApi(url string, inJSON interface{}, outJSON interface{}, w http.ResponseWriter) error {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(inJSON)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("AZURE_FACE_KEY"))
	client := &http.Client{}
	if res, err := client.Do(req); res.StatusCode != 200 {
		return err
	} else {
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&outJSON); err != nil {
			return err
		}
	}
	return nil
}

func FaceDetect(url string, w http.ResponseWriter) (faceDetectResList []FaceDetectResponse, err error) {
	faceDetectURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/detect"
	inJSON := FaceDetectRequest{URL: url}
	err = PostAzureApi(faceDetectURL, inJSON, &faceDetectResList, w)
	return faceDetectResList, err
}

func FaceVerify(faceID string, personID string, w http.ResponseWriter) (faceVerifyRes FaceVerifyResponse, err error) {
	faceVerifyURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/verify"
	inJSON := FaceVerifyRequest{FaceID: faceID, PersonID: personID, PersonGroupID: personGroupID}
	err = PostAzureApi(faceVerifyURL, inJSON, &faceVerifyRes, w)
	return faceVerifyRes, err
}

func CreatePerson(name string, userData string, w http.ResponseWriter) (string, error) {
	createPersonURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/persongroups/" + personGroupID + "/persons"
	inJSON := CreatePersonRequest{Name: name, UserData: userData}
	var createPersonRes CreatePersonResponse
	err := PostAzureApi(createPersonURL, inJSON, &createPersonRes, w)
	return createPersonRes.PersonID, err
}

func FaceIdentification(db *gorm.DB, w http.ResponseWriter, url string) ([]string, error) {
	var allusers models.Users
	allusers.GetAllUsers(db)

	faceDetectResList, err := FaceDetect(url, w)
	if err != nil {
		return nil, err
	}

	l.Debugf("faceDetectResList: %+v\n\n", faceDetectResList)
	downloadUserIDs := make([]string, 0)
	for _, faceDetectRes := range faceDetectResList {
		for _, user := range allusers.Users {
			faceVerifyRes, err := FaceVerify(faceDetectRes.FaceID, user.AzurePersonID, w)
			if err != nil {
				return nil, err
			}
			l.Debugf("faceVerifyRes: %+v\n\n", faceVerifyRes)

			if faceVerifyRes.IsIdentical == true {
				download := models.Download{UserID: user.UserID, URL: url}
				if err := download.CreateRecord(db); err != nil {
					return nil, err
				}
				downloadUserIDs = append(downloadUserIDs, download.UserID)
			}
		}
	}
	return downloadUserIDs, nil
}
