package dashboard

import (
	"context"
	pb_dashboard "stackv2/go/proto/dashboard"
	pb_feature_toggle "stackv2/go/proto/feature_toggle"

	empty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *DashboardServer) CreateFeatureToggle(context.Context, *pb_dashboard.CreateFeatureToggleRequest) (*pb_feature_toggle.FeatureToggle, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFeatureToggle not implemented")
}
func (s *DashboardServer) ListFeatureToggles(context.Context, *pb_dashboard.ListFeatureToggleRequest) (*pb_dashboard.ListFeatureToggleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFeatureToggles not implemented")
}
func (s *DashboardServer) GetFeatureToggle(context.Context, *pb_dashboard.GetFeatureToggleRequest) (*pb_feature_toggle.FeatureToggle, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeatureToggle not implemented")
}
func (s *DashboardServer) GetFeatureToggleHistory(context.Context, *pb_dashboard.GetFeatureToggleHistoryRequest) (*pb_feature_toggle.FeatureToggleHistory, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeatureToggleHistory not implemented")
}
func (s *DashboardServer) UpdateFeatureToggle(context.Context, *pb_dashboard.UpdateFeatureToggleRequest) (*pb_feature_toggle.FeatureToggle, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFeatureToggle not implemented")
}
func (s *DashboardServer) DeleteFeatureToggle(context.Context, *pb_dashboard.DeleteFeatureToggleRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFeatureToggle not implemented")
}
