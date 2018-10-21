package awsutils

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/aws/credentials"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"time"
)

func getS3Client() {
	creds := credentials.NewStaticCredentials("Your AWS Key", "Your AWS Secret", "")
	client := s3.New(creds, s3.BucketLocationConstraintApNortheast1, nil)
	return client
}

func UploadS3(body *File) string {
	client := getS3Client()
	params := &s3.PutObjectInput{
		Bucket: aws.String("ivy-west-winter"), // Required
		Key:    aws.String("your_s3_key"),     // Required
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(body),
	}
	req, err := client.PutObjectRequest(&s3.PutObject(params))
	urlStr, err := req.Presign(10 * time.Minute)
	return urlStr
}
