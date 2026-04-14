//ff:func feature=pkg-auth type=util control=sequence
//ff:what 저장된 해시와 평문 비밀번호가 일치하는지 검증한다
package auth

import "golang.org/x/crypto/bcrypt"

// @func verifyPassword
// @error 401
// @description 저장된 해시와 평문 비밀번호가 일치하는지 검증한다

func VerifyPassword(req VerifyPasswordRequest) (VerifyPasswordResponse, error) {
	err := bcrypt.CompareHashAndPassword([]byte(req.PasswordHash), []byte(req.Password))
	return VerifyPasswordResponse{}, err
}
