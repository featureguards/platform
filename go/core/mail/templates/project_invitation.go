package templates

const (
	projectInvitationPrefix = "project-invite"
)

type ProjectInvitationTemplate struct {
	defaults
	m *ProjectInvite
}

type ProjectInvite struct {
	FirstName string
	Email     string
	Sender    string
	Project   string
	Link      string
}

func NewProjectInvitationTemplate(m *ProjectInvite) *ProjectInvitationTemplate {
	return &ProjectInvitationTemplate{m: m}
}

func (t *ProjectInvitationTemplate) ToName() string {
	return t.m.FirstName
}

func (t *ProjectInvitationTemplate) ToEmail() string {
	return t.m.Email
}

func (t *ProjectInvitationTemplate) Subject() (string, error) {
	return loadTextTemplate(projectInvitationPrefix+".subject.gotmpl", t.m)
}

func (t *ProjectInvitationTemplate) Body() (string, error) {
	return loadTextTemplate(projectInvitationPrefix+".body.gotmpl", t.m)
}

func (t *ProjectInvitationTemplate) HtmlBody() (string, error) {
	return loadTextTemplate(projectInvitationPrefix+".html.body.gotmpl", t.m)
}
