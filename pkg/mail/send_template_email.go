//ff:func feature=pkg-mail type=util control=sequence
//ff:what Go 템플릿으로 HTML 이메일을 발송한다
package mail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strconv"
)

// @func sendTemplateEmail
// @description Go 템플릿으로 HTML 이메일을 발송한다

// SendTemplateEmail accepts ctx for signature parity. See SendEmail comment.
func SendTemplateEmail(ctx context.Context, req SendTemplateEmailRequest) (SendTemplateEmailResponse, error) {
	_ = ctx
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	tmpl, err := template.New("email").Parse(req.TemplateName)
	if err != nil {
		return SendTemplateEmailResponse{}, err
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, nil); err != nil {
		return SendTemplateEmailResponse{}, err
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		from, req.To, req.Subject, body.String())
	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%d", host, port)
	err = smtp.SendMail(addr, auth, from, []string{req.To}, []byte(msg))
	return SendTemplateEmailResponse{}, err
}
