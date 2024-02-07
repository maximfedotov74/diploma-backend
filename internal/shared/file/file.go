package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"github.com/maximfedotov74/diploma-backend/internal/domain/model"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var once sync.Once
var minioClient *minio.Client

type FileClient struct {
	minio      *minio.Client
	mainBucket string
}

func New(minioUrl string, user string, password string, bucket string, ctx context.Context) *FileClient {
	once.Do(func() {
		client, err := minio.New(minioUrl, &minio.Options{Creds: credentials.NewStaticV4(user, password, ""), Secure: false})
		if err != nil {
			log.Fatalf("Failed to connect to minio service, cause: %s", err.Error())
		}

		exists, err := client.BucketExists(ctx, bucket)

		if err != nil {
			log.Fatalf("Failed to start File Client, cause: %s", err.Error())
		}

		if !exists {
			err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: ""})
			if err != nil {
				log.Fatalf("Failed when make new bucket, cause: %s", err.Error())
			}
			policy := `{"Version":"2012-10-17","Statement":[{"Action":["s3:GetObject"],"Effect":"Allow","Principal":"*","Resource":["arn:aws:s3:::` + bucket + `/*"],"Sid":""}]}`

			err = client.SetBucketPolicy(ctx, bucket, policy)

			if err != nil {
				log.Fatalf("Failed when set policy to bucket: %s, cause: %s", bucket, err.Error())
			}
		}

		minioClient = client
	})
	return &FileClient{minio: minioClient, mainBucket: bucket}
}

func (c *FileClient) Upload(ctx context.Context, h *multipart.FileHeader) (*model.UploadResponse, error) {
	file, err := h.Open()
	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("Error when open file, cause: %s", err.Error())
	}

	contentType := h.Header.Get("Content-type")
	fileBytes, err := io.ReadAll(file)

	if err != nil {
		return nil, fmt.Errorf("Error when get file bytes with io: %s", err.Error())
	}

	splittedContentType := strings.Split(contentType, "/")
	fileType := splittedContentType[0]
	extType := splittedContentType[1]
	ext := strings.TrimPrefix(filepath.Ext(h.Filename), filepath.Base(h.Filename))

	fileName := uuid.New().String()

	if fileType == "image" && extType != "svg+xml" {
		compressOptions := bimg.Options{Quality: 50, Type: bimg.WEBP}
		webpBytes, err := bimg.Resize(fileBytes, compressOptions)
		if err != nil {
			return nil, fmt.Errorf("Error when compressing image, cause: %s", err.Error())
		}

		newType := http.DetectContentType(webpBytes)
		newName := fileName + ".webp"
		reader := bytes.NewReader(webpBytes)
		_, err = c.minio.PutObject(ctx, c.mainBucket, newName, reader, reader.Size(), minio.PutObjectOptions{
			ContentType:  newType,
			UserMetadata: map[string]string{"x-amz-acl": "public-read"},
		})
		if err != nil {
			return nil, fmt.Errorf("Error when uploading file, cause: %s", err.Error())
		}

		return &model.UploadResponse{Path: path.Join("/", "storage", c.mainBucket, newName)}, nil
	}

	reader := bytes.NewReader(fileBytes)
	newName := fileName + ext
	_, err = c.minio.PutObject(ctx, c.mainBucket, newName, reader, reader.Size(), minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	})
	if err != nil {
		return nil, fmt.Errorf("Error when uploading file, cause: %s", err.Error())
	}

	return &model.UploadResponse{Path: path.Join("/", "storage", c.mainBucket, newName)}, nil

}
