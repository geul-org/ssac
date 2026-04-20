//ff:func feature=pkg-session type=util control=sequence
//ff:what 세션에서 key를 삭제한다
package session

import "context"

// @func delete
// @description 세션에서 key를 삭제한다

func Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error) {
	return DeleteResponse{}, defaultModel.Delete(ctx, req.Key)
}
