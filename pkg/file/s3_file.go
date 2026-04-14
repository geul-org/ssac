//ff:type feature=pkg-file type=model
//ff:what S3 기반 파일 저장소 구조체
package file

import "github.com/aws/aws-sdk-go-v2/service/s3"

type s3File struct {
	client *s3.Client
	bucket string
}
