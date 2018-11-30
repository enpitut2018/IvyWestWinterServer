package faceidentification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
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

func PostAzureApi(url string, inJSON interface{}, outJSON interface{}, w http.ResponseWriter) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(inJSON)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("AZURE_FACE_KEY"))
	client := &http.Client{}
	if res, _ := client.Do(req); res.StatusCode != 200 {
		httputils.RespondError(w, http.StatusBadRequest, res.Status)
		panic(res.Status)
	} else {
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&outJSON); err != nil {
			httputils.RespondError(w, http.StatusBadRequest, err.Error())
			panic(err.Error())
		}
	}
}

func FaceDetect(url string, w http.ResponseWriter) (faceDetectResList []FaceDetectResponse) {
	faceDetectURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/detect"
	inJSON := FaceDetectRequest{URL: url}
	PostAzureApi(faceDetectURL, inJSON, &faceDetectResList, w)
	return faceDetectResList
}

func FaceVerify(faceID string, personID string, w http.ResponseWriter) (faceVerifyRes FaceVerifyResponse) {
	faceVerifyURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/verify"
	inJSON := FaceVerifyRequest{FaceID: faceID, PersonID: personID, PersonGroupID: personGroupID}
	PostAzureApi(faceVerifyURL, inJSON, &faceVerifyRes, w)
	return faceVerifyRes
}

func CreatePerson(name string, userData string, w http.ResponseWriter) (createPersonRes CreatePersonResponse) {
	createPersonURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/persongroups/" + personGroupID + "/persons"
	inJSON := CreatePersonRequest{Name: name, UserData: userData}
	PostAzureApi(createPersonURL, inJSON, &createPersonRes, w)
	return createPersonRes
}

func FaceIdentification(db *gorm.DB, w http.ResponseWriter, url string) []string {
	var allusers models.Users
	allusers.GetAllUsers(db, w)

	faceDetectResList := FaceDetect(url, w)
	fmt.Printf("faceDetectResList: %+v\n\n", faceDetectResList)

	downloadUserIDs := make([]string, 0)
	for _, faceDetectRes := range faceDetectResList {
		for _, user := range allusers.Users {
			faceVerifyRes := FaceVerify(faceDetectRes.FaceID, user.AzurePersonID, w)
			fmt.Printf("faceVerifyRes: %+v\n\n", faceVerifyRes)
			if faceVerifyRes.IsIdentical == true {
				download := models.Download{UserID: user.UserID, URL: url}
				download.CreateRecord(db, w)
				downloadUserIDs = append(downloadUserIDs, download.UserID)
			}
		}
	}
	return downloadUserIDs
}
