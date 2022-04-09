package users

import (
	"context"
	"stackv2/go/core/app_context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/ory"
	pb_user "stackv2/go/proto/user"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	kratos "github.com/ory/kratos-client-go"
	"gorm.io/gorm"
)

const (
	Me = "me"
)

func FetchUserForSession(ctx context.Context, db *gorm.DB) (*models.User, error) {
	session, ok := app_context.SessionFromContext(ctx)
	if !ok {
		return nil, models.ErrNoSession
	}
	u := &models.User{}
	res := db.WithContext(ctx).First(u, "ory_id", session.Identity.Id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}
	if res.Error != nil {
		err := errors.WithStack(res.Error)
		log.Error(err)
		return nil, err
	}
	return u, nil
}

func FetchIdentityFromUserId(ctx context.Context, userID ids.ID, db *gorm.DB, client *kratos.APIClient) (*kratos.Identity, error) {
	u := &models.User{Model: models.Model{ID: userID}}
	res := db.WithContext(ctx).First(u)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}
	if res.Error != nil {
		err := errors.WithStack(res.Error)
		log.Error(err)
		return nil, err
	}

	return FetchIdentity(ctx, u.OryID, client)

}

func FetchIdentity(ctx context.Context, oryID string, client *kratos.APIClient) (*kratos.Identity, error) {
	req := client.V0alpha2Api.AdminGetIdentity(ctx, oryID)
	identity, _, err := client.V0alpha2Api.AdminGetIdentityExecute(req)
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, err
	}
	return identity, nil
}

func PbUser(identity *kratos.Identity, user *models.User) *pb_user.User {
	traits := ory.Traits(identity.Traits.(map[string]interface{}))
	addresses := make([]*pb_user.User_VerifiableAddress, len(identity.VerifiableAddresses))
	for i, address := range identity.VerifiableAddresses {
		verified := address.Verified
		if !verified {
			if strings.ToLower(traits.Email()) == strings.ToLower(address.Value) {
				verified = traits.EmailVerified()
			}
		}
		addresses[i] = &pb_user.User_VerifiableAddress{
			Address:  address.Value,
			Verified: verified,
		}
	}

	recovery := make([]*pb_user.User_VerifiableAddress, len(identity.RecoveryAddresses))
	for i, address := range identity.RecoveryAddresses {
		recovery[i] = &pb_user.User_VerifiableAddress{
			Address: address.Value,
		}
	}

	return &pb_user.User{
		Id:                string(user.ID),
		OryId:             identity.Id,
		FirstName:         traits.FirstName(),
		LastName:          traits.LastName(),
		Addresses:         addresses,
		RecoveryAddresses: recovery,
		Domain:            traits.Domain(),
		Profile:           traits.Profile(),
	}
}

func FilterPbUser(pb *pb_user.User) {
	pb.Id = ""
	pb.OryId = ""
	pb.RecoveryAddresses = nil
}
