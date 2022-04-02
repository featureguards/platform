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

func (t Traits) FirstName() string {
	return t["first_name"].(string)
}

func (t Traits) LastName() string {
	return t["last_name"].(string)
}

func (t Traits) Domain() string {
	res, ok := t["hd"]
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
	addresses := make([]*pb_dashboard.User_VerifiableAddress, len(session.Identity.VerifiableAddresses))
	for i, address := range session.Identity.VerifiableAddresses {
		addresses[i] = &pb_dashboard.User_VerifiableAddress{
			Address:  address.Value,
			Verified: address.Verified,
		}
	}

	recovery := make([]*pb_dashboard.User_VerifiableAddress, len(session.Identity.RecoveryAddresses))
	for i, address := range session.Identity.RecoveryAddresses {
		recovery[i] = &pb_dashboard.User_VerifiableAddress{
			Address: address.Value,
		}
	}

	return &pb_dashboard.User{
		Id:                string(user.ID),
		OryId:             session.Identity.Id,
		Active:            session.GetActive(),
		FirstName:         traits.FirstName(),
		LastName:          traits.LastName(),
		Addresses:         addresses,
		RecoveryAddresses: recovery,
		Domain:            traits.Domain(),
		Profile:           traits.Profile(),
	}
}
