package mail

type EmailTemplate interface {
	EmailSubject() (string, error)
	EmailBody() (string, error)
	EmailHtmlBody() (string, error)
	EmailRecipientEmail() (string, error)
	EmailRecipientName() (string, error)
	FromName() (string, error)
	FromEmail() (string, error)
}
