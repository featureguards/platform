package mail

import (
	"context"
	"crypto/tls"
	"net/url"
	"strconv"

	gomail "github.com/ory/mail/v3"
	"github.com/pkg/errors"
)

type Opts struct {
	URL *url.URL
}

type Courier struct {
	client *gomail.Dialer
}

func New(opts Opts) (*Courier, error) {
	url := opts.URL
	port, err := strconv.Atoi(url.Port())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sslSkipVerify, _ := strconv.ParseBool(url.Query().Get("skip_ssl_verify"))

	serverName := url.Query().Get("server_name")
	if serverName == "" {
		serverName = url.Hostname()
	}
	password, _ := url.User.Password()
	client := gomail.NewDialer(url.Hostname(), port, url.User.Username(), password)

	var tlsCertificates []tls.Certificate
	// SMTP schemes
	// smtp: smtp clear text (with uri parameter) or with StartTLS (enforced by default)
	// smtps: smtp with implicit TLS (recommended way in 2021 to avoid StartTLS downgrade attacks
	//    and defaulting to fully-encrypted protocols https://datatracker.ietf.org/doc/html/rfc8314)
	switch url.Scheme {
	case "smtp":
		// Enforcing StartTLS by default for security best practices (config review, etc.)
		skipStartTLS, _ := strconv.ParseBool(url.Query().Get("disable_starttls"))
		if !skipStartTLS {
			// #nosec G402 This is ok (and required!) because it is configurable and disabled by default.
			client.TLSConfig = &tls.Config{InsecureSkipVerify: sslSkipVerify, Certificates: tlsCertificates, ServerName: serverName}
		}
	case "smtps":
		// #nosec G402 This is ok (and required!) because it is configurable and disabled by default.
		client.TLSConfig = &tls.Config{InsecureSkipVerify: sslSkipVerify, Certificates: tlsCertificates, ServerName: serverName}
		client.SSL = true
	}

	return &Courier{
		client: client,
	}, nil
}

func (c *Courier) Send(ctx context.Context, t EmailTemplate) error {
	m := gomail.NewMessage()

	m.SetAddressHeader("From", t.FromEmail(), t.FromName())
	m.SetAddressHeader("To", t.ToEmail(), t.ToName())
	subject, err := t.Subject()
	if err != nil {
		return err
	}
	m.SetHeader("Subject", subject)

	htmlBody, err := t.HtmlBody()
	if err != nil {
		return err
	}
	body, err := t.Body()
	if err != nil {
		return err
	}

	m.SetBody("text/plain", body)
	m.AddAlternative("text/html", htmlBody)
	mailer, err := c.client.Dial(ctx)
	if err != nil {
		return err
	}
	defer mailer.Close()
	return gomail.Send(ctx, mailer, m)
}
