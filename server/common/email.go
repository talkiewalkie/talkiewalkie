package common

import (
	"log"
	"net/smtp"
)

type EmailClient interface {
	SendEmail(content []byte, to []string) error
}

var _ EmailClient = GmailClient{}

type GmailClient struct{}

func (c GmailClient) SendEmail(content []byte, to []string) error {
	auth := smtp.PlainAuth("", "accounts@talkiewalkie.app", "", "smtp.gmail.com")
	log.Print("sending email")
	if err := smtp.SendMail("smtp.gmail.com:587", auth, "TalkieWalkie <accounts@talkiewalkie.app>", to, content); err != nil {
		return err
	}
	return nil
}

func initEmailClient() EmailClient {
	return GmailClient{}
}
