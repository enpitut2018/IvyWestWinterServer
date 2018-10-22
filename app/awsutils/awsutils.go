package awsutils

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	ACCESS_KEY = ""
	SECRET_KEY = ""
)

func getS3Client() *s3.S3 {
	creds := credentials.NewStaticCredentials(ACCESS_KEY, SECRET_KEY, "")
	client := s3.New(session.New(), &aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-northeast-1"),
	})
	return client
}

func UploadS3(body []byte, keyname string) error {
	client := getS3Client()
	input := &s3.PutObjectInput{
		Bucket: aws.String("ivy-west-winter"), // Required
		Key:    aws.String(keyname),           // Required
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(body),
	}
	_, err := client.PutObject(input)
	return err
}
