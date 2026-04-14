//ff:func feature=pkg-storage type=util control=sequence
//ff:what S3 클라이언트를 생성한다
package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func newS3Client(endpoint, region string) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	opts := func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		}
	}
	return s3.NewFromConfig(cfg, opts), nil
}
