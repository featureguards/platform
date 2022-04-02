package templates

import (
	"path/filepath"
)

const (
	projectInvitationPrefix = "project-invitation-subject"
)

type (
	ProjectInvitationTemplate struct {
		m *ProjectInvitation
	}
	ProjectInvitation struct {
		To  string
		URL string
	}
)

func NewProjectInvitationTemplate(m *ProjectInvitation) *ProjectInvitationTemplate {
	return &ProjectInvitationTemplate{m: m}
}

func (t *ProjectInvitationTemplate) EmailRecipient() (string, error) {
	return t.m.To, nil
}

func (t *ProjectInvitationTemplate) EmailSubject() (string, error) {
	return loadTextTemplate(filepath.Join(path, projectInvitationPrefix+"-subject"), t.m)
}

func (t *ProjectInvitationTemplate) EmailBody() (string, error) {
	return loadTextTemplate(filepath.Join(path, projectInvitationPrefix+"-body"), t.m)
}
