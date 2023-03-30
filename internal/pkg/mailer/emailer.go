package mailer

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/whitewolf185/mangaparser/internal/config"
	"gopkg.in/gomail.v2"
)

//go:generate mockgen -destination=./mock/email_getter_mock.go -package=mock github.com/whitewolf185/mangaparser/internal/pkg/mailer EmailGetter
type EmailGetter interface {
	GetEmailByID(ctx context.Context, personID uuid.UUID) (string, error)
}

type EbookMailer struct{
	smtpServer string
	smtpPort int
	adressGetter EmailGetter
}

func NewEbookMailer(adressGetter EmailGetter) EbookMailer {
	return EbookMailer{
		smtpServer: "smtp.mail.ru",
		smtpPort: 465,
		adressGetter: adressGetter,
	}
}

func (em EbookMailer) SendManga(ctx context.Context, personID uuid.UUID, mangaFilePath string) error {
	toEmailAdress, err := em.adressGetter.GetEmailByID(ctx, personID)
	if err != nil{
		return errors.Wrap(err, "failred to get email adress from db")
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