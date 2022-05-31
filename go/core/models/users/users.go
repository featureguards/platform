package users

import (
	"context"
	"platform/go/core/app_context"
	"platform/go/core/ids"
	"platform/go/core/models"
	"platform/go/core/ory"

	pb_user "github.com/featureguards/client-go/proto/user"

	"github.com/pkg/errors"

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
	return GetByOryID(ctx, session.Identity.Id, db)
}

func GetByOryID(ctx context.Context, oryID string, db *gorm.DB) (*models.User, error) {
	u := &models.User{}
	res := db.WithContext(ctx).First(u, "ory_id", oryID)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}
	if res.Error != nil {
		return nil, errors.WithStack(res.Error)
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
		return nil, errors.WithStack(res.Error)
	}

	return FetchIdentity(ctx, u.OryID, client)

}

func FetchIdentity(ctx context.Context, oryID string, client *kratos.APIClient) (*kratos.Identity, error) {
	req := client.V0alpha2Api.AdminGetIdentity(ctx, oryID)
	identity, _, err := client.V0alpha2Api.AdminGetIdentityExecute(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return identity, nil
}

func Pb(identity *kratos.Identity, user *models.User) *pb_user.User {
	traits := ory.Traits(identity.Traits.(map[string]interface{}))
	addresses := make([]*pb_user.User_VerifiableAddress, len(identity.VerifiableAddresses))
	for i, address := range identity.VerifiableAddresses {
		addresses[i] = &pb_user.User_VerifiableAddress{
			Address:  address.Value,
			Verified: address.Verified,
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
	pb.OryId = ""
	pb.RecoveryAddresses = nil
}

func LimitedPbUser(pb *pb_user.User) {
	FilterPbUser(pb)
	pb.Addresses = nil
	pb.Profile = ""

}
