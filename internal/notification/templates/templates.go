package templates

type MailTemplate string

const (
	REGISTRATION MailTemplate = "registration-verification-template.html"
	RESET_PASSWORD_TOKEN MailTemplate = "change-password-template.html"
)

