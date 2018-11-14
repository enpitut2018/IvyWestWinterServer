package awsutils

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

var (
	ACCESS_KEY = os.Getenv("AWS_ACCESS_KEY")
	SECRET_KEY = os.Getenv("AWS_SECRET_KEY")
)

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
