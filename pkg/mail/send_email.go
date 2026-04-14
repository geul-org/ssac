//ff:func feature=pkg-mail type=util control=sequence
//ff:what SMTP를 통해 이메일을 발송한다
package mail

import (
	"fmt"
	"net/smtp"
)

// @func sendEmail
// @description SMTP를 통해 이메일을 발송한다

func SendEmail(req SendEmailRequest) (SendEmailResponse, error) {
	auth := smtp.PlainAuth("", req.Username, req.Password, req.Host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		req.From, req.To, req.Subject, req.Body)
	addr := fmt.Sprintf("%s:%d", req.Host, req.Port)
	err := smtp.SendMail(addr, auth, req.From, []string{req.To}, []byte(msg))
	return SendEmailResponse{}, err
}
