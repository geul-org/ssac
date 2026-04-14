//ff:type feature=pkg-mail type=model
//ff:what 템플릿 이메일 발송 요청 모델
package mail

type SendTemplateEmailRequest struct {
	To           string
	Subject      string
	TemplateName string
}
