//ff:func feature=pkg-storage type=util control=sequence
//ff:what S3 호환 스토리지에 파일을 업로드한다
package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// @func uploadFile
// @description S3 호환 스토리지에 파일을 업로드한다

func UploadFile(ctx context.Context, req UploadFileRequest) (UploadFileResponse, error) {
	client, err := newS3Client(ctx, req.Endpoint, req.Region)
	if err != nil {
		return UploadFileResponse{}, err
	}
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(req.Bucket),
		Key:         aws.String(req.Key),
		Body:        bytes.NewReader(req.Data),
		ContentType: aws.String(req.ContentType),
	})
	if err != nil {
		return UploadFileResponse{}, err
	}
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", req.Bucket, req.Region, req.Key)
	if req.Endpoint != "" {
		url = fmt.Sprintf("%s/%s/%s", req.Endpoint, req.Bucket, req.Key)
	}
	return UploadFileResponse{URL: url}, nil
}
