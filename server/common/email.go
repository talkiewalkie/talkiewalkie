package common

import (
	"log"
	"net/smtp"
)

type EmailClient struct{}

func (c EmailClient) SendEmail(to string, content []byte) error {
	auth := smtp.PlainAuth("", "accounts@talkiewalkie.app", "", "smtp.gmail.com")
	log.Print("sending email")
	if err := smtp.SendMail("smtp.gmail.com:587", auth, "TalkieWalkie <accounts@talkiewalkie.app>", []string{to}, content); err != nil {
		return err
	}
	return nil
}

func initEmailClient() EmailClient {
	return EmailClient{}
}
