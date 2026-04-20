//ff:func feature=pkg-mail type=util control=sequence
//ff:what SMTP를 통해 이메일을 발송한다
package mail

import (
	"context"
	"fmt"
	"net/smtp"
)

// @func sendEmail
// @description SMTP를 통해 이메일을 발송한다

// SendEmail accepts ctx as the first argument for request-cancellation
// propagation. net/smtp itself does not support ctx yet; the parameter is
// held for future migration to ctx-aware SMTP dialers and signature parity
// with the rest of the ssac/pkg runtime.
func SendEmail(ctx context.Context, req SendEmailRequest) (SendEmailResponse, error) {
	_ = ctx
	auth := smtp.PlainAuth("", req.Username, req.Password, req.Host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		req.From, req.To, req.Subject, req.Body)
	addr := fmt.Sprintf("%s:%d", req.Host, req.Port)
	err := smtp.SendMail(addr, auth, req.From, []string{req.To}, []byte(msg))
	return SendEmailResponse{}, err
}
