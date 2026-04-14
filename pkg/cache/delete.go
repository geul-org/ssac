//ff:func feature=pkg-cache type=util control=sequence
//ff:what 캐시에서 key를 삭제한다
package cache

import "context"

// @func delete
// @description 캐시에서 key를 삭제한다

func Delete(req DeleteRequest) (DeleteResponse, error) {
	return DeleteResponse{}, defaultModel.Delete(context.Background(), req.Key)
}
