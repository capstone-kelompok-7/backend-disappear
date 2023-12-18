package email

import (
	"errors"
	"github.com/wneessen/go-mail"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type EmailSenderInterface interface {
	EmailService(email, otp string) error
}

type Sender struct{}

func NewEmailService() EmailSenderInterface {
	return &Sender{}
}

func (s *Sender) EmailService(email, otp string) error {
	secretUser := os.Getenv("SMTP_USER")
	secretPass := os.Getenv("SMTP_PASS")
	secretPort := os.Getenv("SMTP_PORT")

	convPort, err := strconv.Atoi(secretPort)
	if err != nil {
		return err
	}

	m := mail.NewMsg()
	if err := m.From(secretUser); err != nil {
		return err
	}
	if err := m.To(email); err != nil {
		return err
	}

	m.Subject("Verifikasi Email - Disappear Organization")
	emailTemplate := struct {
		OTP   string
		Email string
	}{
		OTP:   otp,
		Email: email,
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("Failed to get the current file path")
	}

	templatePath := filepath.Join(filepath.Dir(filename), "template.html")

	tmpl, err := template.New("emailTemplate").ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var bodyContent strings.Builder
	if err := tmpl.ExecuteTemplate(&bodyContent, "template.html", emailTemplate); err != nil {
		return err
	}

	m.SetBodyString(mail.TypeTextHTML, bodyContent.String())

	c, err := mail.NewClient("smtp.gmail.com", mail.WithPort(convPort), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(secretUser), mail.WithPassword(secretPass))
	if err != nil {
		return err
	}
	if err := c.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
