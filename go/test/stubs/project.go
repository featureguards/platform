package stubs

import (
	"context"

	pb_dashboard "platform/go/proto/dashboard"

	"github.com/Pallinder/go-randomdata"
)

func (s *Stubs) createProject(ctx context.Context) error {
	name := randomdata.Alphanumeric(10)
	envs := []string{randomdata.RandStringRunes(6), randomdata.RandStringRunes(6)}
	proj, err := s.App.DashboardClient.CreateProject(s.WithToken(ctx), &pb_dashboard.CreateProjectRequest{
		Name:         name,
		Environments: []*pb_dashboard.CreateProjectRequest_NewEnvironment{{Name: envs[0]}, {Name: envs[1]}},
	})
	if err != nil {
		return err
	}

	s.Proj = proj
	return nil
}
