package faceidentification

import (
	"fmt"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"bytes"
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
)

type FaceDetectRequest struct {
	URL string `json:"url"`
}

type FaceDetectResponse struct {
	FaceID string
	FaceRectangle string
	FaceAttributes string
}

type FaceVerifyRequest struct {
	FaceID string `json:"faceid"`
	PersonID string `json:"personid"`
}

type FaceVerifyResponse struct {
	IsIdentical bool
	Confidence float32
}

func PostAzureApi(url string, inJSON interface{}, outJSON interface{}, w http.ResponseWriter) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(inJSON)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("AZURE_FACE_KEY"))
	client := &http.Client{}
	fmt.Printf("%+v\n", req)
    if res, err := client.Do(req); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	} else {
		decoder := json.NewDecoder(res.Body)
		if err = decoder.Decode(&outJSON); err != nil {
			httputils.RespondError(w, http.StatusBadRequest, err.Error())
			panic(err.Error())
		}
	}
	fmt.Printf("%+v\n", outJSON)
} 

func FaceDetect(url string, w http.ResponseWriter) (faceDetectResList []FaceDetectResponse) {
	faceDetectURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/detect"
	inJSON := FaceDetectRequest{URL: url}
	var outJSON []FaceDetectResponse
	PostAzureApi(faceDetectURL, inJSON, outJSON, w)
	return faceDetectResList
}

func FaceVerify(faceID string, personID string, w http.ResponseWriter) (faceVerifyRes FaceVerifyResponse) {
	faceVerifyURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/verify"
	inJSON := FaceVerifyRequest{FaceID: faceID, PersonID: personID}
	var outJSON FaceVerifyResponse
	PostAzureApi(faceVerifyURL, inJSON, outJSON, w)
	return faceVerifyRes
}

func FaceIdentification(db *gorm.DB, w http.ResponseWriter, url string) {
	var allusers models.Users
	allusers.GetAllUsers(db, w)
	fmt.Printf("%+v\n", allusers)

	faceDetectResList := FaceDetect(url, w)
	fmt.Printf("%+v\n", faceDetectResList)

	for _, faceDetectRes := range faceDetectResList {
		for _, user := range allusers.Users {
			// faceVerifyRes := FaceVerify(faceDetectRes.FaceID, user.AzurePersonID)
			faceVerifyRes := FaceVerify(faceDetectRes.FaceID, "0b4bbd63-ff70-423b-9aff-5263c745ff98", w)
			fmt.Printf("%+v\n", faceVerifyRes)
			if faceVerifyRes.IsIdentical == true {
				download := models.Download{UserID: user.UserID, URL: url}
				download.CreateRecord(db, w)
			}
		}
	}
}