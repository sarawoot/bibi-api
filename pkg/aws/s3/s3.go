package s3

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	bucket   string
	s3Client *s3.Client
}

// New:
// required: env as follow
// 	AWS_REGION
// 	AWS_ACCESS_KEY_ID
// 	AWS_SECRET_ACCESS_KEY
func New(bucket string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Client{
		bucket:   bucket,
		s3Client: s3.NewFromConfig(cfg),
	}, nil
}

func (c *Client) UploadFile(ctx context.Context, key string, uploadFile io.Reader) error {
	uploader := manager.NewUploader(c.s3Client)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, uploadFile); err != nil {
		return err
	}

	contentType := http.DetectContentType(buf.Bytes()[:512])

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        buf,
		ContentType: &contentType,
	})

	return err
}

// func getFileContentType(out io.Reader) (string, error) {
// 	buf := bytes.NewBuffer(nil)
// 	if _, err := io.Copy(buf, out); err != nil {
// 		return "", err
// 	}

// 	// Use the net/http package's handy DectectContentType function. Always returns a valid
// 	// content-type by returning "application/octet-stream" if no others seemed to match.
// 	contentType := http.DetectContentType(buf.Bytes()[:512])

// 	return contentType, nil
// }
