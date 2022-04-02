package mail

import (
	"context"
	"os"

	"github.com/pkg/errors"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Opts struct {
	URL string
}

type Courier struct {
	URL    string
	client *sendgrid.Client
}

func New(opts Opts) (*Courier, error) {
	url := opts.URL
	if url == "" {
		url = os.Getenv("SENDGRID_API_KEY")
	}
	if url == "" {
		return nil, errors.Errorf("no URL set")
	}

	client := sendgrid.NewSendClient(url)
	return &Courier{
		client: client,
	}, nil
}

func (m *Courier) Send(ctx context.Context, t EmailTemplate) error {
	body, err := t.EmailBody()
	if err != nil {
		return err
	}

	htmlBody, err := t.EmailHtmlBody()
	if err != nil {
		return err
	}

	subject, err := t.EmailSubject()
	if err != nil {
		return err
	}

	recipient, err := t.EmailRecipientName()
	if err != nil {
		return err
	}
	recipientEmail, err := t.EmailRecipientEmail()
	if err != nil {
		return err
	}
	sender, err := t.FromName()
	if err != nil {
		return err
	}
	senderEmail, err := t.FromEmail()
	if err != nil {
		return err
	}

	from := mail.NewEmail(sender, senderEmail)
	to := mail.NewEmail(recipient, recipientEmail)
	message := mail.NewSingleEmail(from, subject, to, body, htmlBody)
	if _, err := m.client.SendWithContext(ctx, message); err != nil {
		return err
	}
	return nil
}
