//ff:type feature=pkg-mail type=model
//ff:what 이메일 발송 요청 모델
package mail

type SendEmailRequest struct {
	Host     string // SMTP 호스트 (예: smtp.gmail.com)
	Port     int    // SMTP 포트 (예: 587)
	Username string
	Password string
	From     string
	To       string
	Subject  string
	Body     string // plain text
}
