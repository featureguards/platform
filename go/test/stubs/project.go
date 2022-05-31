package stubs

import (
	"context"

	pb_dashboard "platform/go/proto/dashboard"

	"github.com/Pallinder/go-randomdata"
)

func (s *Stubs) createProject(ctx context.Context) error {
	name := randomdata.Noun()
	envs := []string{randomdata.Noun(), randomdata.Noun()}
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
