package model

// @dto
// Refund는 환불 계산 결과다 (DDL 테이블 없음).
type Refund struct {
	RefundAmount int64
	RefundRate   float64
}
