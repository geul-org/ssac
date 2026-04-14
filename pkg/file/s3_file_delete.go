//ff:func feature=pkg-file type=util control=sequence
//ff:what S3 Delete — DeleteObject로 객체 삭제
package file

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (f *s3File) Delete(ctx context.Context, key string) error {
	_, err := f.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(key),
	})
	return err
}
