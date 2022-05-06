package dashboard

import (
	"context"
	"fmt"

	"stackv2/go/core/app"
	pb_dashboard "stackv2/go/proto/dashboard"

	"gorm.io/gorm"
)

// Make sure DashboardServer implements everything
var _ pb_dashboard.DashboardServer = &DashboardServer{}

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

func (s *DashboardServer) DB(ctx context.Context) *gorm.DB {
	return s.app.DB.WithContext(ctx)
}
