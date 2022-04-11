package dashboard

import (
	"context"
	"stackv2/go/core/app_context"
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/users"
	pb_dashboard "stackv2/go/proto/dashboard"
	pb_user "stackv2/go/proto/user"

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
	u, err := users.FetchUserForSession(ctx, s.app.DB)
	if errors.Is(err, models.ErrNotFound) {
		// Create the object
		userID, err := ids.IDFromShard(ids.ShardIDFromKey(session.Identity.Id), ids.User)
		if err != nil {
			err := errors.WithStack(err)
			log.Error(err)
			return nil, status.Errorf(codes.Internal, "could not retrive user")
		}
		u.ID = userID
		u.OryID = session.Identity.Id
		res := s.app.DB.WithContext(ctx).FirstOrCreate(u)
		if res.Error != nil {
			err := errors.WithStack(res.Error)
			log.Error(err)
			return nil, status.Errorf(codes.Internal, "could not retrive user")
		}
	}

	return users.PbUser(&session.Identity, u), nil
}
