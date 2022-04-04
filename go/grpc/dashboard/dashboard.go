package dashboard

import (
	"fmt"

	"stackv2/go/core/app"
	pb_dashboard "stackv2/go/proto/dashboard"
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
