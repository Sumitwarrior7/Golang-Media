package mailer

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile string, username string, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(fromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// Template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// The subject buffer will contain the static string "Finish Registration with Golang Media"
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	// Resending emails if some error occurs till MaxRetries, but if any try email is sent successfully then this function stops
	for i := 0; i < MaxRetries; i++ {
		response, retryErr := m.client.Send(message)
		if retryErr != nil {
			// Exponential Backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email sent with status code : %v", response.StatusCode)
		return nil
	}

	return fmt.Errorf("Failed to send email after %d attempts", MaxRetries)
}
