package email

import (
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"github.com/JooseMM/nutripia-backend-api/internal/notification/templates"
	"github.com/JooseMM/nutripia-backend-api/pkg/core"
)

type MailerSender interface {
	Send() *core.BaseError
}

type GmailSender struct {
	Port     string
	Host     string
	Password string
	From     string
	To       []string

	Subject string
	Message []byte
}

func NewMailSender(
	to []string,
	subject string,
	template templates.MailTemplate,
	replacements map[string]string,
) (MailerSender, *core.BaseError) {
	port := os.Getenv("SMTP_PORT")
	host := os.Getenv("SMTP_HOST")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_EMAIL")

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	prepare, prepareErr := prepareTemplate(template, replacements)
	if prepareErr != nil {
		return nil, prepareErr
	}

	return &GmailSender{
		Port:     port,
		Host:     host,
		Password: password,
		From:     from,

		Subject: subject,
		Message: []byte("Subject: " + subject + "\n" + mime + *prepare),
		To:      to,
	}, nil
}

func (s *GmailSender) Send() *core.BaseError {
	auth := smtp.PlainAuth("", s.From, s.Password, s.Host)

	err := smtp.SendMail(s.Host+":"+s.Port, auth, s.From, s.To, s.Message)
	if err != nil {
		return core.UnexpectedError(err.Error())
	}
	return nil
}

func prepareTemplate(
	templateName templates.MailTemplate,
	replacements map[string]string,
) (*string, *core.BaseError) {
	path := getAbsPath(string(templateName))

	content, contentErr := os.ReadFile(*path)
	if contentErr != nil {
		return nil, core.UnexpectedError(contentErr.Error())
	}

	var oldNew []string
	for k, v := range replacements {
		oldNew = append(oldNew, k, v)
	}

	replacer := strings.NewReplacer(oldNew...)
	final := replacer.Replace(string(content))
	return &final, nil
}

func getAbsPath(templateFile string) *string {
	path := filepath.Join("internal", "notifications", "templates", templateFile)
	return &path
}
