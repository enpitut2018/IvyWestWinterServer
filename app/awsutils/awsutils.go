package awsutils

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
	"encoding/base64"
	"path/filepath"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/rs/xid"
)

var (
	ACCESS_KEY = os.Getenv("AWS_ACCESS_KEY")
	SECRET_KEY = os.Getenv("AWS_SECRET_KEY")
)

func UploadPhoto(w http.ResponseWriter, base64Str string, s3FolderPath string) string {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic("can't decode base64")
	}
	guid := xid.New() // xidというユニークなID
	imageFileName := guid.String()+".jpg"
	if false == UploadS3(data, filepath.Join(s3FolderPath ,imageFileName)) {
		httputils.RespondError(w, http.StatusInternalServerError, "Can't upload the photo.")
		panic("Can't upload the photo.")
	}
	return filepath.Join("https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter", s3FolderPath, imageFileName)
}

func getS3Client() *s3.S3 {
	creds := credentials.NewStaticCredentials(ACCESS_KEY, SECRET_KEY, "")
	client := s3.New(session.New(), &aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-northeast-1"),
	})
	return client
}

func UploadS3(body []byte, keyname string) bool {
	client := getS3Client()
	input := &s3.PutObjectInput{
		Bucket: aws.String("ivy-west-winter"), // Required
		Key:    aws.String(keyname),           // Required
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(body),
	}
	if _, err := client.PutObject(input); err != nil {
		panic(err.Error())
		return false
	} else {
		return true
	}
}
