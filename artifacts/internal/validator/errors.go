package validator

import "fmt"

// ValidationError는 검증 에러 하나를 나타낸다.
type ValidationError struct {
	FileName string // 원본 파일명
	FuncName string // 함수명
	SeqIndex int    // sequence 인덱스
	Tag      string // 관련 태그 (e.g. "@model", "@action")
	Message  string // 에러 메시지
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s:%s:seq[%d] %s — %s", e.FileName, e.FuncName, e.SeqIndex, e.Tag, e.Message)
}
