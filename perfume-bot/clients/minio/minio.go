package minio

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type Client struct {
	mcl      *minio.Client
	bucket   string
	endpoint string
}

func New(mcl *minio.Client, bucket string, endpoint string) *Client {
	return &Client{
		mcl:      mcl,
		bucket:   bucket,
		endpoint: endpoint,
	}
}

func (c *Client) UploadPhoto(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("file.Open: %w", err)
	}
	info, err := c.mcl.PutObject(
		ctx,
		c.bucket,
		file.Filename,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", fmt.Errorf("mcl.PutObject: %w", err)
	}
	return fmt.Sprintf("http://%s/%s/%s", c.endpoint, c.bucket, info.Key), nil
}
