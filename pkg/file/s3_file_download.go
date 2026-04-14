//ff:func feature=pkg-file type=util control=sequence
//ff:what S3 Download — GetObject로 객체 다운로드
package file

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (f *s3File) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := f.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}
