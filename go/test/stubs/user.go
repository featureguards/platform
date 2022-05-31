package stubs

import (
	"context"

	pb_dashboard "platform/go/proto/dashboard"
)

func (s *Stubs) createUser(ctx context.Context) error {
	_, token, err := s.App.CreateUserWithSession(ctx)
	if err != nil {
		return err
	}

	s.Token = token
	// This lazily creates the user object
	user, err := s.App.DashboardClient.GetUser(s.WithToken(ctx), &pb_dashboard.GetUserRequest{UserId: "me"})
	if err != nil {
		return err
	}
	s.User = user
	return nil
}
