package notification

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/fabiokaelin/f-oauth/config"
)

func SendEmail(toAddress string, subject string, body string) error {
	host := config.EmailSMTPHost
	port := config.EmailSMTPPort
	from := config.EmailFromAddress
	name := config.EmailFromName
	password := config.EmailSMTPPassword

	if host == "" || port == "" || from == "" || password == "" {
		return fmt.Errorf("email config is incomplete: EMAIL_SMTP_HOST, EMAIL_SMTP_PORT, EMAIL_FROM_ADDRESS and EMAIL_SMTP_PASSWORD must be set")
	}

	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", name, from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", toAddress))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := host + ":" + port
	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(addr, auth, from, []string{toAddress}, []byte(msg.String()))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func SendPasswordResetEmail(toAddress string, resetLink string) error {
	subject := "Reset your password"
	body := fmt.Sprintf(
		"Hi,\n\nYou requested a password reset for your %s account.\n\nClick the link below to set a new password:\n%s\n\nThis link is valid for 3 hours. If you did not request a password reset, you can safely ignore this email.\n\nRegards,\n%s",
		config.EmailFromName, resetLink, config.EmailFromName,
	)

	err := SendEmail(toAddress, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	return nil
}
