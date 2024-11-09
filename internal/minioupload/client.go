package minioupload

import (
	"context"
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
)

func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (MinioApi, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

func (m *Client) CreateBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	exists, err := m.ExistBucket(ctx, bucketName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("bucket already exists")
	}
	return m.client.MakeBucket(ctx, bucketName, opts)
}

func (m *Client) RemoveBucket(ctx context.Context, bucketName string) error {
	return m.client.RemoveBucket(ctx, bucketName)
}

func (m *Client) ExistBucket(ctx context.Context, bucketName string) (bool, error) {
	exists, err := m.client.BucketExists(ctx, bucketName)
	return exists, err
}

func (m *Client) SetTagBucket(ctx context.Context, bucketName string, tags *tags.Tags) error {
	return m.client.SetBucketTagging(ctx, bucketName, tags)
}

func (m *Client) RemoveTagBucket(ctx context.Context, bucketName string) error {
	return m.client.RemoveBucketTagging(ctx, bucketName)
}

func (m *Client) UploadFile(ctx context.Context, bucketName, objectName string, file *os.File, size int64, opt minio.PutObjectOptions) (*minio.UploadInfo, error) {
	uploadInfo, err := m.client.PutObject(ctx, bucketName, objectName, file, size, opt)
	if err != nil {
		return nil, err
	}
	return &uploadInfo, nil
}

func (m *Client) DownloadUrlFile(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objectName)
	presignedUrl, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expires, reqParams)
	if err != nil {
		return nil, err
	}
	return presignedUrl, nil
}

func (m *Client) RemoveFile(ctx context.Context, bucketName, objectName string, opt minio.RemoveObjectOptions) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, opt)
}
