//ff:func feature=pkg-file type=util control=sequence
//ff:what 저장소에서 파일을 삭제한다
package file

import "context"

// @func delete
// @description 저장소에서 파일을 삭제한다

func Delete(req DeleteRequest) (DeleteResponse, error) {
	return DeleteResponse{}, defaultModel.Delete(context.Background(), req.Key)
}
