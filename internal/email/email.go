// Package email allows sending transactional emails
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
