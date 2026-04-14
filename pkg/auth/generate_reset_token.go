//ff:func feature=pkg-auth type=util control=sequence
//ff:what 비밀번호 리셋용 32바이트 랜덤 hex 토큰을 생성한다
package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// @func generateResetToken
// @description 비밀번호 리셋용 32바이트 랜덤 hex 토큰을 생성한다

func GenerateResetToken(req GenerateResetTokenRequest) (GenerateResetTokenResponse, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return GenerateResetTokenResponse{}, err
	}
	return GenerateResetTokenResponse{Token: hex.EncodeToString(b)}, nil
}
