package dashboard

import (
	"context"
	"platform/go/core/app_context"
	"platform/go/core/ids"
	"platform/go/core/models"
	"platform/go/core/models/users"
	pb_dashboard "platform/go/proto/dashboard"

	pb_user "github.com/featureguards/featureguards-go/v2/proto/user"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DashboardServer) GetUser(ctx context.Context, req *pb_dashboard.GetUserRequest) (*pb_user.User, error) {
	if req.UserId != users.Me {
		// TODO: Support admin queries
		return nil, status.Error(codes.NotFound, "no user found")
	}
	session, ok := app_context.SessionFromContext(ctx)
	if !ok {
		// This should never happen since we enforce this at the middleware.
		log.Fatal("received unauthenticated session.")
		return nil, status.Error(codes.Unauthenticated, "invalid session")
	}
	// Make sure that we have a user. Otherwise, create the user if a new user
	u, err := users.FetchUserForSession(ctx, s.app.DB())
	if errors.Is(err, models.ErrNotFound) {
		// Create the object
		userID, err := ids.IDFromShard(ids.ShardIDFromKey(session.Identity.Id), ids.User)
		if err != nil {
			log.Errorf("%s\n", err)
			return nil, status.Errorf(codes.Internal, "could not retrive user")
		}
		u = &models.User{
			Model: models.Model{ID: userID},
			OryID: session.Identity.Id,
		}
		res := s.DB(ctx).FirstOrCreate(u)
		if res.Error != nil {
			log.Errorf("%s\n", errors.WithStack(res.Error))
			return nil, status.Errorf(codes.Internal, "could not retrive user")
		}
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return users.Pb(&session.Identity, u), nil
}
