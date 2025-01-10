package usecase

import (
	"context"
	"fmt"
	"mime/multipart"

	// "ps-beli-mang/configs"
	// "ps-beli-mang/pkg/errs"
	// "ps-beli-mang/pkg/helper"

	"github.com/SkandarEverest/refresh-golang/internal/exception"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type FileUseCase struct {
	Client *s3.Client
	Bucket string
	Config *viper.Viper
}

func NewFileUseCase(logger *logrus.Logger, cfg *viper.Viper) *FileUseCase {
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.GetString("S3_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.GetString("S3_ID"), cfg.GetString("S3_SECRET_KEY"), "")),
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(awsConfig)

	return &FileUseCase{
		Client: s3Client,
		Bucket: cfg.GetString("S3_BUCKET_NAME"),
		Config: cfg,
	}
}

func (c *FileUseCase) UploadFile(file multipart.File) (*string, error) {
	defer file.Close()

	fileName := uuid.New().String() + ".jpeg"

	_, err := c.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(fileName),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   file,
	})
	if err != nil {
		return nil, exception.ServerError(err.Error())
	}

	fileUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.Bucket, c.Config.GetString("S3_REGION"), fileName)
	return &fileUrl, nil
}
