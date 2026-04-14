//ff:func feature=pkg-auth type=util control=sequence
//ff:what 평문 비밀번호를 bcrypt 해시로 변환한다
package auth

import "golang.org/x/crypto/bcrypt"

// @func hashPassword
// @description 평문 비밀번호를 bcrypt 해시로 변환한다

func HashPassword(req HashPasswordRequest) (HashPasswordResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	return HashPasswordResponse{HashedPassword: string(hash)}, err
}
