package dashboard

import (
	"context"
	"errors"
	"fmt"

	"stackv2/go/core/app"
	"stackv2/go/core/app_context"
	id "stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/models/users"
	pb_dashboard "stackv2/go/proto/dashboard"

	log "github.com/sirupsen/logrus"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// server is used to implement dashboard.DashboardServer.
type DashboardServer struct {
	pb_dashboard.UnimplementedDashboardServer
	app *app.App
}

type DashboardOpts struct {
	App *app.App
}

func New(opts DashboardOpts) (*DashboardServer, error) {
	if opts.App == nil {
		return nil, fmt.Errorf("app must be specified")
	}
	return &DashboardServer{app: opts.App}, nil
}

func (s *DashboardServer) Me(ctx context.Context, req *pb_dashboard.MeRequest) (*pb_dashboard.User, error) {
	session, ok := app_context.SessionFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid session")
	}
	// Make sure that we have a user. Otherwise, create the user if a new user
	u := &models.User{}
	res := s.app.DB.WithContext(ctx).First(u, "ory_id", session.Identity.Id)
	log.Warnf("%+v", res)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// Create the object
		userID, err := s.app.IDs.IDFromShard(s.app.IDs.ShardIDFromKey(session.Identity.Id), id.User)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s", err)
		}
		u.ID = userID
		u.OryID = session.Identity.Id
		res = s.app.DB.WithContext(ctx).FirstOrCreate(u)
		if res.Error != nil {
			return nil, err
		}
	}

	return users.PbUser(session, u), nil
}
func (s *DashboardServer) CreateProject(context.Context, *pb_dashboard.CreateProjectRequest) (*pb_dashboard.Project, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProject not implemented")
}
func (s *DashboardServer) ListProjects(context.Context, *pb_dashboard.ListProjectsRequest) (*pb_dashboard.ListProjectsResponse, error) {
	res := &pb_dashboard.ListProjectsResponse{
		Projects: []*pb_dashboard.Project{{Id: "111", Name: "Foo"}, {Id: "222", Name: "Bar"}},
	}
	return res, nil
}
func (s *DashboardServer) GetProject(context.Context, *pb_dashboard.GetProjectRequest) (*pb_dashboard.Project, error) {
	return &pb_dashboard.Project{Id: "111", Name: "Foo"}, nil
}
func (s *DashboardServer) DeleteProject(context.Context, *pb_dashboard.DeleteProjectRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProject not implemented")
}
