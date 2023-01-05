package providers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Store struct {
	s3         *minio.Client
	logger     *zap.Logger
	BucketName string
	domain     string
}

func NewStoreProvider(vip *viper.Viper, logger *zap.Logger) (*Store, error) {
	if !vip.IsSet("s3") {
		return nil, errors.New("")
	}
	config := vip.Sub("s3")

	endpoint := config.GetString("endpoint")
	accessKeyID := config.GetString("secretId")
	secretAccessKey := config.GetString("secretKey")
	useSSL := config.GetBool("useSSL")
	buckerName := config.GetString("buckerName")
	domain := config.GetString("domain")
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}
	return &Store{s3: minioClient, BucketName: buckerName, domain: domain, logger: logger}, nil
}

func (s *Store) PresignedUrl(ctx context.Context, path, filename string) string {

	if strings.HasPrefix(path, s.domain) {
		path = strings.TrimPrefix(path, s.domain)
	}

	if strings.HasPrefix(path, "http") || path == "" {
		return path
	}

	reqParams := make(url.Values)
	if filename != "" {
		reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	presignedURL, err := s.s3.PresignedGetObject(ctx, s.BucketName, path, time.Second*24*60*60, reqParams)
	if err != nil {
		panic(err)
	}
	return presignedURL.String()
}

func (s *Store) PutObject(ctx context.Context, path string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (string, error) {
	oss, err := s.s3.PutObject(ctx, s.BucketName, path, reader, objectSize, opts)
	if err != nil {
		return "", err
	}
	s.logger.Info("Successfully uploaded", zap.String("objectName", path), zap.Any("info", oss))
	s.logger.Info("PutObject", zap.String("Location", oss.Location))
	return fmt.Sprintf("%s/%s", s.domain, path), nil
}

func (s *Store) DelObject(ctx context.Context, filepath string) (string, error) {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	filepath = strings.TrimPrefix(filepath, "/")
	err := s.s3.RemoveObject(ctx, s.BucketName, filepath, opts)
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("%s/%s", s.domain, filepath), err
	}
	s.logger.Info("Successfully deleted", zap.String("objectName", filepath))
	return fmt.Sprintf("%s/%s", s.domain, filepath), nil
}

func (s *Store) GetPolicyToken(ctx context.Context, filepath string) (*url.URL, error) {
	return s.s3.PresignedPutObject(ctx, s.BucketName, filepath, time.Hour)
}
