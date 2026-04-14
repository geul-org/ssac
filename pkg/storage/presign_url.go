//ff:func feature=pkg-storage type=util control=sequence
//ff:what 서명된 다운로드 URL을 생성한다
package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// @func presignURL
// @description 서명된 다운로드 URL을 생성한다

func PresignURL(req PresignURLRequest) (PresignURLResponse, error) {
	client, err := newS3Client(req.Endpoint, req.Region)
	if err != nil {
		return PresignURLResponse{}, err
	}
	expiresIn := req.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 3600
	}
	presigner := s3.NewPresignClient(client)
	presigned, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(req.Bucket),
		Key:    aws.String(req.Key),
	}, s3.WithPresignExpires(time.Duration(expiresIn)*time.Second))
	if err != nil {
		return PresignURLResponse{}, err
	}
	return PresignURLResponse{URL: presigned.URL}, nil
}
