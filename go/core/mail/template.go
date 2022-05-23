package mail

type EmailTemplate interface {
	Subject() (string, error)
	Body() (string, error)
	HtmlBody() (string, error)
	ToEmail() string
	ToName() string
	FromName() string
	FromEmail() string
}
