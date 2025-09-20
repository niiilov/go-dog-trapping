package s3

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	creds "github.com/niiilov/go-dog-trapping/internal/config"
)

type Storage struct {
	client        *s3.Client
	defautlBucket string
}

func New(cred *creds.AwsCreds) *Storage {

	cfg := aws.Config{
		Region: cred.Region,
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(cred.AccessKeyID, cred.SecretAccessKey, ""),
		),
		BaseEndpoint: &cred.BaseEndpoint,
	}

	// Создаем клиента для доступа к хранилищу S3
	client := s3.NewFromConfig(cfg)

	// Запрашиваем список всех файлов в бакете
	result, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("sobaki"),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range result.Contents {
		log.Printf("object=%s size=%d Bytes last modified=%s", aws.ToString(object.Key), aws.ToInt64(object.Size), object.LastModified.Local().Format("2006-01-02 15:04:05 Monday"))
	}

	return &Storage{
		client:        client,
		defautlBucket: "sobaki",
	}

}

func (s *Storage) UploadFile(ctx context.Context, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.defautlBucket),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
				log.Printf("Error while uploading object to %s. The object is too large.\n"+
					"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
					"or the multipart upload API (5TB max).", s.defautlBucket)
			} else {
				log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
					fileName, s.defautlBucket, objectKey, err)
			}
		} else {
			err = s3.NewObjectExistsWaiter(s.client).Wait(
				ctx, &s3.HeadObjectInput{Bucket: aws.String(s.defautlBucket), Key: aws.String(objectKey)}, time.Minute)
			if err != nil {
				log.Printf("Failed attempt to wait for object %s to exist.\n", objectKey)
			}
		}
	}
	return err
}

func (s *Storage) GetFileURL(objectKey string) string {

	fmt.Println("я тут")
	presignClient := s3.NewPresignClient(s.client)
	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.defautlBucket),
		Key:    aws.String(objectKey),
	}, s3.WithPresignExpires(15*time.Minute))

	fmt.Println(presignResult.URL)
	if err != nil {
		log.Printf("Couldn't get presigned URL for %v/%v. Here's why: %v\n", s.defautlBucket, objectKey, err)
		return ""
	}
	return presignResult.URL
}
