// Package main just to test Brevo API client
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading %s file", ".env")
	}

	var ctx context.Context
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", os.Getenv("BREVO_API_KEY"))
	cfg.AddDefaultHeader("partner-key", os.Getenv("BREVO_API_KEY"))

	br := brevo.NewAPIClient(cfg)

	result, resp, err := br.AccountApi.GetAccount(ctx)
	if err != nil {
		fmt.Println("Error when calling AccountApi->get_account: ", err.Error())
		return
	}

	fmt.Println("GetAccount Object:", result, " GetAccount Response: ", resp)

	_, _, err = br.TransactionalEmailsApi.SendTransacEmail(context.TODO(), brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  os.Getenv("EMAIL_FROM_NAME"),
			Email: os.Getenv("EMAIL_FROM_ADDRESS"),
		},
		To: []brevo.SendSmtpEmailTo{
			{
				Email: os.Getenv("EMAIL_TO_ADDRESS"),
				Name:  os.Getenv("EMAIL_TO_NAME"),
			},
		},
		Subject:     "Test Email from Brevo API",
		TextContent: "This is a test email sent using the Brevo API.",
	})
	if err != nil {
		fmt.Println("Error when calling TransactionalEmailsApi->sendTransacEmail: ", err.Error())
		return
	}
}
