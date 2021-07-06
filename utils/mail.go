package utils

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/AnhHoangQuach/go-intern-spores/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//Defind Requets for send mail...
type MailSender struct {
	from    string
	to      []string
	subject string
	body    string
}

// Create request to send mail content...
func NewMailSender(to []string, subject string) *MailSender {
	return &MailSender{
		to:      to,
		subject: subject,
	}
}

func (r *MailSender) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *MailSender) sendMail() bool {
	mailConfig := config.GetMailOption()

	from := mail.NewEmail("Welups", mailConfig.MailFrom)
	subject := r.subject
	to := mail.NewEmail("", r.to[0])
	htmlContent := r.body
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(mailConfig.SendGridApiKey)
	_, err := client.Send(message)

	if err != nil {
		return false
	} else {
		return true
	}
}

func (r *MailSender) Send(templateName string, items interface{}) error {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		return err
	}
	if ok := r.sendMail(); ok {
		return nil
	} else {
		return fmt.Errorf("Failed to send the mail to %s\n", r.to)
	}
}
