package util

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"sync"

	"github.com/spf13/viper"
)

type SMTPClient struct {
	From   string
	Client *smtp.Client
	Lock   sync.Mutex 
}

func NewSMTP(viper *viper.Viper) *SMTPClient {
	host := viper.GetString("email.host")
	sender := viper.GetString("email.sender")
	password := viper.GetString("email.password")
	port := viper.GetString("email.port")

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), &tls.Config{
		ServerName: host,
	})
	if err != nil {
		log.Fatal(err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Fatal(err)
	}

	auth := smtp.PlainAuth("", sender, password, host)
	if err := client.Auth(auth); err != nil {
		log.Fatal(err)
	}

	return &SMTPClient{
		From:   sender,
		Client: client,
	}
}

func (s *SMTPClient) SendMail(subject, body string, to ...string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	if err := s.Client.Mail(s.From); err != nil {
		return err
	}

	for _, recipient := range to {
		if err := s.Client.Rcpt(recipient); err != nil {
			return err
		}
	}

	w, err := s.Client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *SMTPClient) Close() {
	s.Client.Quit()
}

