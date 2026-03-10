package billing

import (
	"time"
)

// @func calculateRefund
// @description 예약 취소 시 환불 금액을 계산한다 (시작일까지 남은 기간 기반)

type CalculateRefundRequest struct {
	ReservationID int64
	StartAt       time.Time
	EndAt         time.Time
}

type CalculateRefundResponse struct {
	RefundAmount int64
	RefundRate   float64
}

func CalculateRefund(req CalculateRefundRequest) (CalculateRefundResponse, error) {
	now := time.Now()
	daysUntilStart := time.Until(req.StartAt).Hours() / 24

	var rate float64
	switch {
	case daysUntilStart >= 7:
		rate = 1.0 // 전액 환불
	case daysUntilStart >= 3:
		rate = 0.5 // 50% 환불
	case daysUntilStart >= 1:
		rate = 0.3 // 30% 환불
	default:
		rate = 0.0 // 환불 불가
	}

	totalDays := req.EndAt.Sub(req.StartAt).Hours() / 24
	dailyRate := int64(10000) // 1일 기본 요금 (실제론 Room 가격 참조)
	totalAmount := dailyRate * int64(totalDays)
	refundAmount := int64(float64(totalAmount) * rate)

	_ = now // suppress unused

	return CalculateRefundResponse{
		RefundAmount: refundAmount,
		RefundRate:   rate,
	}, nil
}
