package users

import (
	"stackv2/go/core/models"
	pb_dashboard "stackv2/go/proto/dashboard"

	kratos "github.com/ory/kratos-client-go"
)

type Traits map[string]interface{}

func (t Traits) Email() string {
	return t["email"].(string)
}

func (t Traits) EmailVerified() bool {
	return t["email_verified"].(bool)
}

func (t Traits) FirstName() string {
	return t["first_name"].(string)
}

func (t Traits) LastName() string {
	return t["last_name"].(string)
}

func (t Traits) Domain() string {
	res, ok := t["domain"]
	if ok {
		return res.(string)
	}
	return ""
}

func (t Traits) Profile() string {
	res, ok := t["profile"]
	if ok {
		return res.(string)
	}
	return ""
}

func PbUser(session *kratos.Session, user *models.User) *pb_dashboard.User {
	traits := Traits(session.Identity.Traits.(map[string]interface{}))
	return &pb_dashboard.User{
		Id:                string(user.ID),
		OryId:             session.Identity.Id,
		Active:            session.GetActive(),
		FirstName:         traits.FirstName(),
		LastName:          traits.LastName(),
		Addresses:         []*pb_dashboard.User_VerifiableAddress{{Address: traits.Email(), Verified: traits.EmailVerified()}},
		RecoveryAddresses: []*pb_dashboard.User_VerifiableAddress{{Address: traits.Email(), Verified: traits.EmailVerified()}},
		Domain:            traits.Domain(),
		Profile:           traits.Profile(),
	}
}
