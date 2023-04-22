package mailer

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"

	"github.com/whitewolf185/mangaparser/internal/config"
	customerrors "github.com/whitewolf185/mangaparser/pkg/custom_errors"
)

type EbookMailer struct{
	smtpServer string
	smtpPort int
}

func NewEbookMailer() EbookMailer {
	return EbookMailer{
		smtpServer: "smtp.mail.ru",
		smtpPort: 465,
	}
}

func (em EbookMailer) SendManga(ctx context.Context, toEmailAdress string, mangaFilePath string) error {
	if toEmailAdress == "" {
		return customerrors.ErrEmailsNotFound
	}

	message := gomail.NewMessage()
	from := config.GetValue(config.EmailAccount)
	message.SetHeader("From", from) 
	message.SetHeader("To", toEmailAdress)
	message.Attach(mangaFilePath)

	password := config.GetValue(config.EmailPassword)
	dialer := gomail.NewDialer(em.smtpServer, em.smtpPort, from, password)

	if err := dialer.DialAndSend(message); err != nil{
		return errors.Wrap(err, "cannot send email with manga")
	}
	return nil
}