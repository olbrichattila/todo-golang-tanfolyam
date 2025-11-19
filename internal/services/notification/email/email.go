package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	repositoryContracts "todo/internal/repositories/contracts"
	"todo/internal/services/contracts"
)

type EmailData struct {
	Token  string
	AppURL string
}

func New(appConfig repositoryContracts.AppConfig) contracts.NotificationService {
	// TODO nil check
	return &notification{
		appConfig: appConfig,
	}
}

type notification struct {
	appConfig repositoryContracts.AppConfig
}

func (n *notification) SendRegistrationNotification(email string, token string) error {
	tmpl, err := template.ParseFiles("templates/email/register.html")
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth("", n.appConfig.SmtpUsername(), n.appConfig.SmtpPassword(), n.appConfig.SmtpHost())

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, EmailData{Token: token, AppURL: n.appConfig.AppURL()}); err != nil {
		return err
	}

	body := buf.String()

	// MIME headers
	subject := "Confirm Your Registration"
	msg := fmt.Sprintf("From: %s\r\n", n.appConfig.EmailFrom())
	msg += fmt.Sprintf("To: %s\r\n", email)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"
	msg += body

	addr := fmt.Sprintf("%s:%s", n.appConfig.SmtpHost(), n.appConfig.SmtpPort())
	return smtp.SendMail(addr, auth, n.appConfig.EmailFrom(), []string{email}, []byte(msg))
}
