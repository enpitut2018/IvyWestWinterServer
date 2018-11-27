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
	FaceRectangle interface{}
	FaceAttributes interface{}
}

type FaceVerifyRequest struct {
	FaceID string `json:"faceId"`
	PersonID string `json:"personId"`
	PersonGroupID string `json:"personGroupId"`
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
	fmt.Printf("req: %+v\n\n", req)
    if res, _ := client.Do(req); res.StatusCode != 200 {
		fmt.Printf("res: %+v\n\n", res)
		httputils.RespondError(w, http.StatusBadRequest, res.Status)
		// panic(res.Status)
	} else {
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&outJSON); err != nil {
			httputils.RespondError(w, http.StatusBadRequest, err.Error())
			// panic(err.Error())
		}
	}
} 

func FaceDetect(url string, w http.ResponseWriter) (faceDetectResList []FaceDetectResponse) {
	faceDetectURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/detect"
	//url = "https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/upload-photos/fukuyama.jpeg"
	inJSON := FaceDetectRequest{URL: url}
	PostAzureApi(faceDetectURL, inJSON, &faceDetectResList, w)
	return faceDetectResList
}

func FaceVerify(faceID string, personID string, w http.ResponseWriter) (faceVerifyRes FaceVerifyResponse) {
	faceVerifyURL := "https://japaneast.api.cognitive.microsoft.com/face/v1.0/verify"
	inJSON := FaceVerifyRequest{FaceID: faceID, PersonID: personID, PersonGroupID: "ivy-west-winter-test"}
	PostAzureApi(faceVerifyURL, inJSON, &faceVerifyRes, w)
	return faceVerifyRes
}

func FaceIdentification(db *gorm.DB, w http.ResponseWriter, url string) {
	fmt.Printf("url: %+v\n\n", url)
	var allusers models.Users
	allusers.GetAllUsers(db, w)
	fmt.Printf("allusers: %+v\n\n", allusers)

	faceDetectResList := FaceDetect(url, w)
	fmt.Printf("faceDetectResList: %+v\n\n", faceDetectResList)
	if len(faceDetectResList) == 0 {
		httputils.RespondError(w, http.StatusBadRequest, "Can't Detect Face.")
		panic("Can't Detect Face.")
	}

	for _, faceDetectRes := range faceDetectResList {
		for _, user := range allusers.Users {
			// faceVerifyRes := FaceVerify(faceDetectRes.FaceID, user.AzurePersonID)
			faceVerifyRes := FaceVerify(faceDetectRes.FaceID, "0b4bbd63-ff70-423b-9aff-5263c745ff98", w)
			fmt.Printf("faceVerifyRes: %+v\n\n", faceVerifyRes)
			if faceVerifyRes.IsIdentical == true {
				download := models.Download{UserID: user.UserID, URL: url}
				download.CreateRecord(db, w)
			}
		}
	}
}