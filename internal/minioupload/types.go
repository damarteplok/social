package minioupload

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/tags"
)

type MinioApi interface {
	CreateBucket(context.Context, string, minio.MakeBucketOptions) error
	RemoveBucket(context.Context, string) error
	ExistBucket(context.Context, string) (bool, error)
	SetTagBucket(context.Context, string, *tags.Tags) error
	RemoveTagBucket(context.Context, string) error
	UploadFile(ctx context.Context, bucketName, objectName string, file *os.File, size int64, opt minio.PutObjectOptions) (*minio.UploadInfo, error)
	DownloadUrlFile(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error)
	RemoveFile(ctx context.Context, bucketName, objectName string, opt minio.RemoveObjectOptions) error
	UploadBpmnOrForm(ctx context.Context, file *os.File, fileName string) (*minio.UploadInfo, error)
	GetObject(ctx context.Context, bucketName, objectName string) (*os.File, error)
}

type Client struct {
	client *minio.Client
}
