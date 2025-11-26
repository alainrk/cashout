// Package email allows sending transactional emails.
// TODO: Implement an interface to be injected where needed.
package email

import (
	"context"

	brevo "github.com/getbrevo/brevo-go/lib"
)

type EmailService struct {
	fromName  string
	fromEmail string
	client    *brevo.APIClient
}

func NewEmailService(apiKey string, fromName string, fromEmail string) (*EmailService, error) {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	cfg.AddDefaultHeader("partner-key", apiKey)

	br := brevo.NewAPIClient(cfg)

	// Check if the account exists and is verified
	_, _, err := br.AccountApi.GetAccount(context.Background())
	if err != nil {
		return nil, err
	}

	return &EmailService{
		fromName:  fromName,
		fromEmail: fromEmail,
		client:    br,
	}, nil
}

func (e *EmailService) SendTransacEmail(toEmail string, subject string, textContent string) error {
	_, _, err := e.client.TransactionalEmailsApi.SendTransacEmail(context.Background(), brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  e.fromName,
			Email: e.fromEmail,
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: toEmail,
				Name:  "",
			},
		},
		Subject:     subject,
		TextContent: textContent,
	})

	return err
}
