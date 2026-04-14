//ff:func feature=pkg-file type=util control=sequence
//ff:what S3 Upload — PutObject로 객체 업로드
package file

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (f *s3File) Upload(ctx context.Context, key string, body io.Reader) error {
	_, err := f.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}
