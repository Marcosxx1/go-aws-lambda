// Package uploaderservice provides functionality for uploading images to an S3 bucket.
package uploaderservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// UploaderAdapter represents a service for uploading images to an S3 bucket.
type UploaderAdapter struct {
	S3Client *s3.Client
}

// NewUploaderAdapter creates a new UploaderAdapter instance configured with the AWS S3 client.
// It returns a pointer to the UploaderAdapter or an error if the AWS configuration fails.
func NewUploaderAdapter() (*UploaderAdapter, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	return &UploaderAdapter{S3Client: s3Client}, nil
}

// UploadImage uploads the given image to the S3 bucket.
// It takes the image bytes, tabloid ID, and order as parameters.
// It returns the key under which the image is stored in the S3 bucket or an error if upload fails.
func (adapter *UploaderAdapter) UploadImage(image []byte, tabloidID int64, order int) (string, error) {
	if image == nil {
		return "", nil
	}

	if err := adapter.validateImage(image); err != nil {
		return "", err
	}

	key := adapter.getImageKey(image, tabloidID, order)

	_, err := adapter.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_S3_BUCKET_NAME_S3")),
		Key:         aws.String(key),
		ContentType: aws.String(http.DetectContentType(image)),
		Body:        bytes.NewReader(image),
	})
	if err != nil {
		fmt.Println(err)
		return "", errors.New("ERROR_UPLOAD_IMAGE")
	}

	return key, nil
}

// getImageKey generates a unique key for the image based on tabloid ID, order, and UUID.
func (adapter *UploaderAdapter) getImageKey(image []byte, tabloidID int64, order int) string {
	extension := adapter.getImageExtension(image)
	pagina := order + 1
	randomId := uuid.New()
	return fmt.Sprintf("RPA/v3/%d/campanha-%d-%s-pagina-%d%s", tabloidID, tabloidID, randomId, pagina, extension)
}

// validateImage checks if the image has a valid content type.
// It returns an error if the image type is not supported.
func (adapter *UploaderAdapter) validateImage(image []byte) error {
	contentType := http.DetectContentType(image)
	if contentType != "image/png" && contentType != "image/jpg" && contentType != "image/jpeg" {
		return errors.New("invalid image type")
	}
	return nil
}

// getImageExtension extracts the file extension from the image content type.
func (adapter *UploaderAdapter) getImageExtension(image []byte) string {
	contentType := http.DetectContentType(image)
	parts := strings.Split(contentType, "/")
	if len(parts) < 2 {
		return ""
	}
	return "." + parts[1]
}

/* //Utilizando singleton:
package uploaderservice

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type UploaderAdapter struct {
	S3Client *s3.Client
}

var uploader *UploaderAdapter

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader = &UploaderAdapter{S3Client: s3Client}
}

func GetUploaderAdapter() *UploaderAdapter {
	return uploader
}

func (adapter *UploaderAdapter) UploadImage(image []byte, tabloidID int64, order int) (string, error) {
	if image == nil {
		return "", nil
	}

	if err := adapter.validateImage(image); err != nil {
		return "", err
	}

	key := adapter.getImageKey(image, tabloidID, order)

	_, err := adapter.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("AWS_S3_BUCKET_NAME_S3")),
		Key:         aws.String(key),
		ContentType: aws.String(http.DetectContentType(image)),
		Body:        bytes.NewReader(image),
	})
	if err != nil {
		return "", errors.New("ERROR_UPLOAD_IMAGE") //TODO melhorar tratamento de erros
	}

	return key, nil
}

func (adapter *UploaderAdapter) getImageKey(image []byte, tabloidID int64, order int) string {
	extension := adapter.getImageExtension(image)
	pagina := order + 1
	uuid := uuid.New()
	return fmt.Sprintf("RPA/v3/%d/campanha-%d-%s-pagina-%d%s", tabloidID, tabloidID, uuid, pagina, extension)
}

func (adapter *UploaderAdapter) validateImage(image []byte) error {
	contentType := http.DetectContentType(image)
	if contentType != "image/png" && contentType != "image/jpg" && contentType != "image/jpeg" {
		return errors.New("invalid image type") //TODO melhorar tratamento de erros
	}
	return nil
}

func (adapter *UploaderAdapter) getImageExtension(image []byte) string {
	contentType := http.DetectContentType(image)
	parts := strings.Split(contentType, "/")
	if len(parts) < 2 {
		return ""
	}
	return "." + parts[1]
}


*/
