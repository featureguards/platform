// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package dashboard

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	feature_toggle "stackv2/go/proto/feature_toggle"
	project "stackv2/go/proto/project"
	user "stackv2/go/proto/user"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DashboardClient is the client API for Dashboard service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DashboardClient interface {
	// Users
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*user.User, error)
	// Projects
	CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*project.Project, error)
	ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error)
	GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*project.Project, error)
	DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Invites
	CreateProjectInvite(ctx context.Context, in *ProjectInviteRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	ListProjectInvites(ctx context.Context, in *ListProjectInvitesRequest, opts ...grpc.CallOption) (*project.ProjectInvites, error)
	ListUserInvites(ctx context.Context, in *ListUserInvitesRequest, opts ...grpc.CallOption) (*project.ProjectInvites, error)
	ListProjectMembers(ctx context.Context, in *ListProjectMembersRequest, opts ...grpc.CallOption) (*project.ProjectMembers, error)
	GetProjectInvite(ctx context.Context, in *GetProjectInviteRequest, opts ...grpc.CallOption) (*project.ProjectInvite, error)
	UpdateProjectInvite(ctx context.Context, in *UpdateProjectInviteRequest, opts ...grpc.CallOption) (*project.ProjectInvite, error)
	// Environments
	CreateEnvironment(ctx context.Context, in *CreateEnvironmentRequest, opts ...grpc.CallOption) (*project.Environment, error)
	ListEnvironments(ctx context.Context, in *ListEnvironmentsRequest, opts ...grpc.CallOption) (*ListEnvironmentsResponse, error)
	GetEnvironment(ctx context.Context, in *GetEnvironmentRequest, opts ...grpc.CallOption) (*project.Environment, error)
	DeleteEnvironment(ctx context.Context, in *DeleteEnvironmentRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// FeatureToggles
	CreateFeatureToggle(ctx context.Context, in *CreateFeatureToggleRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	ListFeatureToggles(ctx context.Context, in *ListFeatureToggleRequest, opts ...grpc.CallOption) (*ListFeatureToggleResponse, error)
	GetFeatureToggle(ctx context.Context, in *GetFeatureToggleRequest, opts ...grpc.CallOption) (*EnvironmentFeatureToggles, error)
	GetFeatureToggleHistoryForEnvironment(ctx context.Context, in *GetFeatureToggleHistoryRequest, opts ...grpc.CallOption) (*feature_toggle.FeatureToggleHistory, error)
	UpdateFeatureToggle(ctx context.Context, in *UpdateFeatureToggleRequest, opts ...grpc.CallOption) (*feature_toggle.FeatureToggle, error)
	DeleteFeatureToggle(ctx context.Context, in *DeleteFeatureToggleRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type dashboardClient struct {
	cc grpc.ClientConnInterface
}

func NewDashboardClient(cc grpc.ClientConnInterface) DashboardClient {
	return &dashboardClient{cc}
}

func (c *dashboardClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*user.User, error) {
	out := new(user.User)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*project.Project, error) {
	out := new(project.Project)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/CreateProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error) {
	out := new(ListProjectsResponse)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/ListProjects", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*project.Project, error) {
	out := new(project.Project)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/GetProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/DeleteProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) CreateProjectInvite(ctx context.Context, in *ProjectInviteRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/CreateProjectInvite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) ListProjectInvites(ctx context.Context, in *ListProjectInvitesRequest, opts ...grpc.CallOption) (*project.ProjectInvites, error) {
	out := new(project.ProjectInvites)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/ListProjectInvites", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) ListUserInvites(ctx context.Context, in *ListUserInvitesRequest, opts ...grpc.CallOption) (*project.ProjectInvites, error) {
	out := new(project.ProjectInvites)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/ListUserInvites", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) ListProjectMembers(ctx context.Context, in *ListProjectMembersRequest, opts ...grpc.CallOption) (*project.ProjectMembers, error) {
	out := new(project.ProjectMembers)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/ListProjectMembers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) GetProjectInvite(ctx context.Context, in *GetProjectInviteRequest, opts ...grpc.CallOption) (*project.ProjectInvite, error) {
	out := new(project.ProjectInvite)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/GetProjectInvite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) UpdateProjectInvite(ctx context.Context, in *UpdateProjectInviteRequest, opts ...grpc.CallOption) (*project.ProjectInvite, error) {
	out := new(project.ProjectInvite)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/UpdateProjectInvite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) CreateEnvironment(ctx context.Context, in *CreateEnvironmentRequest, opts ...grpc.CallOption) (*project.Environment, error) {
	out := new(project.Environment)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/CreateEnvironment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) ListEnvironments(ctx context.Context, in *ListEnvironmentsRequest, opts ...grpc.CallOption) (*ListEnvironmentsResponse, error) {
	out := new(ListEnvironmentsResponse)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/ListEnvironments", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) GetEnvironment(ctx context.Context, in *GetEnvironmentRequest, opts ...grpc.CallOption) (*project.Environment, error) {
	out := new(project.Environment)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/GetEnvironment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) DeleteEnvironment(ctx context.Context, in *DeleteEnvironmentRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/DeleteEnvironment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) CreateFeatureToggle(ctx context.Context, in *CreateFeatureToggleRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/CreateFeatureToggle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) ListFeatureToggles(ctx context.Context, in *ListFeatureToggleRequest, opts ...grpc.CallOption) (*ListFeatureToggleResponse, error) {
	out := new(ListFeatureToggleResponse)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/ListFeatureToggles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) GetFeatureToggle(ctx context.Context, in *GetFeatureToggleRequest, opts ...grpc.CallOption) (*EnvironmentFeatureToggles, error) {
	out := new(EnvironmentFeatureToggles)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/GetFeatureToggle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) GetFeatureToggleHistoryForEnvironment(ctx context.Context, in *GetFeatureToggleHistoryRequest, opts ...grpc.CallOption) (*feature_toggle.FeatureToggleHistory, error) {
	out := new(feature_toggle.FeatureToggleHistory)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/GetFeatureToggleHistoryForEnvironment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) UpdateFeatureToggle(ctx context.Context, in *UpdateFeatureToggleRequest, opts ...grpc.CallOption) (*feature_toggle.FeatureToggle, error) {
	out := new(feature_toggle.FeatureToggle)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/UpdateFeatureToggle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dashboardClient) DeleteFeatureToggle(ctx context.Context, in *DeleteFeatureToggleRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/dashboard.Dashboard/DeleteFeatureToggle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DashboardServer is the server API for Dashboard service.
// All implementations must embed UnimplementedDashboardServer
// for forward compatibility
type DashboardServer interface {
	// Users
	GetUser(context.Context, *GetUserRequest) (*user.User, error)
	// Projects
	CreateProject(context.Context, *CreateProjectRequest) (*project.Project, error)
	ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error)
	GetProject(context.Context, *GetProjectRequest) (*project.Project, error)
	DeleteProject(context.Context, *DeleteProjectRequest) (*empty.Empty, error)
	// Invites
	CreateProjectInvite(context.Context, *ProjectInviteRequest) (*empty.Empty, error)
	ListProjectInvites(context.Context, *ListProjectInvitesRequest) (*project.ProjectInvites, error)
	ListUserInvites(context.Context, *ListUserInvitesRequest) (*project.ProjectInvites, error)
	ListProjectMembers(context.Context, *ListProjectMembersRequest) (*project.ProjectMembers, error)
	GetProjectInvite(context.Context, *GetProjectInviteRequest) (*project.ProjectInvite, error)
	UpdateProjectInvite(context.Context, *UpdateProjectInviteRequest) (*project.ProjectInvite, error)
	// Environments
	CreateEnvironment(context.Context, *CreateEnvironmentRequest) (*project.Environment, error)
	ListEnvironments(context.Context, *ListEnvironmentsRequest) (*ListEnvironmentsResponse, error)
	GetEnvironment(context.Context, *GetEnvironmentRequest) (*project.Environment, error)
	DeleteEnvironment(context.Context, *DeleteEnvironmentRequest) (*empty.Empty, error)
	// FeatureToggles
	CreateFeatureToggle(context.Context, *CreateFeatureToggleRequest) (*empty.Empty, error)
	ListFeatureToggles(context.Context, *ListFeatureToggleRequest) (*ListFeatureToggleResponse, error)
	GetFeatureToggle(context.Context, *GetFeatureToggleRequest) (*EnvironmentFeatureToggles, error)
	GetFeatureToggleHistoryForEnvironment(context.Context, *GetFeatureToggleHistoryRequest) (*feature_toggle.FeatureToggleHistory, error)
	UpdateFeatureToggle(context.Context, *UpdateFeatureToggleRequest) (*feature_toggle.FeatureToggle, error)
	DeleteFeatureToggle(context.Context, *DeleteFeatureToggleRequest) (*empty.Empty, error)
	mustEmbedUnimplementedDashboardServer()
}

// UnimplementedDashboardServer must be embedded to have forward compatible implementations.
type UnimplementedDashboardServer struct {
}

func (UnimplementedDashboardServer) GetUser(context.Context, *GetUserRequest) (*user.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedDashboardServer) CreateProject(context.Context, *CreateProjectRequest) (*project.Project, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProject not implemented")
}
func (UnimplementedDashboardServer) ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjects not implemented")
}
func (UnimplementedDashboardServer) GetProject(context.Context, *GetProjectRequest) (*project.Project, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProject not implemented")
}
func (UnimplementedDashboardServer) DeleteProject(context.Context, *DeleteProjectRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProject not implemented")
}
func (UnimplementedDashboardServer) CreateProjectInvite(context.Context, *ProjectInviteRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProjectInvite not implemented")
}
func (UnimplementedDashboardServer) ListProjectInvites(context.Context, *ListProjectInvitesRequest) (*project.ProjectInvites, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjectInvites not implemented")
}
func (UnimplementedDashboardServer) ListUserInvites(context.Context, *ListUserInvitesRequest) (*project.ProjectInvites, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUserInvites not implemented")
}
func (UnimplementedDashboardServer) ListProjectMembers(context.Context, *ListProjectMembersRequest) (*project.ProjectMembers, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjectMembers not implemented")
}
func (UnimplementedDashboardServer) GetProjectInvite(context.Context, *GetProjectInviteRequest) (*project.ProjectInvite, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectInvite not implemented")
}
func (UnimplementedDashboardServer) UpdateProjectInvite(context.Context, *UpdateProjectInviteRequest) (*project.ProjectInvite, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProjectInvite not implemented")
}
func (UnimplementedDashboardServer) CreateEnvironment(context.Context, *CreateEnvironmentRequest) (*project.Environment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEnvironment not implemented")
}
func (UnimplementedDashboardServer) ListEnvironments(context.Context, *ListEnvironmentsRequest) (*ListEnvironmentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEnvironments not implemented")
}
func (UnimplementedDashboardServer) GetEnvironment(context.Context, *GetEnvironmentRequest) (*project.Environment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnvironment not implemented")
}
func (UnimplementedDashboardServer) DeleteEnvironment(context.Context, *DeleteEnvironmentRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEnvironment not implemented")
}
func (UnimplementedDashboardServer) CreateFeatureToggle(context.Context, *CreateFeatureToggleRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFeatureToggle not implemented")
}
func (UnimplementedDashboardServer) ListFeatureToggles(context.Context, *ListFeatureToggleRequest) (*ListFeatureToggleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFeatureToggles not implemented")
}
func (UnimplementedDashboardServer) GetFeatureToggle(context.Context, *GetFeatureToggleRequest) (*EnvironmentFeatureToggles, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeatureToggle not implemented")
}
func (UnimplementedDashboardServer) GetFeatureToggleHistoryForEnvironment(context.Context, *GetFeatureToggleHistoryRequest) (*feature_toggle.FeatureToggleHistory, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeatureToggleHistoryForEnvironment not implemented")
}
func (UnimplementedDashboardServer) UpdateFeatureToggle(context.Context, *UpdateFeatureToggleRequest) (*feature_toggle.FeatureToggle, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFeatureToggle not implemented")
}
func (UnimplementedDashboardServer) DeleteFeatureToggle(context.Context, *DeleteFeatureToggleRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFeatureToggle not implemented")
}
func (UnimplementedDashboardServer) mustEmbedUnimplementedDashboardServer() {}

// UnsafeDashboardServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DashboardServer will
// result in compilation errors.
type UnsafeDashboardServer interface {
	mustEmbedUnimplementedDashboardServer()
}

func RegisterDashboardServer(s grpc.ServiceRegistrar, srv DashboardServer) {
	s.RegisterService(&Dashboard_ServiceDesc, srv)
}

func _Dashboard_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).CreateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/CreateProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).CreateProject(ctx, req.(*CreateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_ListProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).ListProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/ListProjects",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).ListProjects(ctx, req.(*ListProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_GetProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).GetProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/GetProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).GetProject(ctx, req.(*GetProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).DeleteProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/DeleteProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).DeleteProject(ctx, req.(*DeleteProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_CreateProjectInvite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectInviteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).CreateProjectInvite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/CreateProjectInvite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).CreateProjectInvite(ctx, req.(*ProjectInviteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_ListProjectInvites_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProjectInvitesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).ListProjectInvites(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/ListProjectInvites",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).ListProjectInvites(ctx, req.(*ListProjectInvitesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_ListUserInvites_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserInvitesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).ListUserInvites(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/ListUserInvites",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).ListUserInvites(ctx, req.(*ListUserInvitesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_ListProjectMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProjectMembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).ListProjectMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/ListProjectMembers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).ListProjectMembers(ctx, req.(*ListProjectMembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_GetProjectInvite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectInviteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).GetProjectInvite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/GetProjectInvite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).GetProjectInvite(ctx, req.(*GetProjectInviteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_UpdateProjectInvite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProjectInviteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).UpdateProjectInvite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/UpdateProjectInvite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).UpdateProjectInvite(ctx, req.(*UpdateProjectInviteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_CreateEnvironment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEnvironmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).CreateEnvironment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/CreateEnvironment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).CreateEnvironment(ctx, req.(*CreateEnvironmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_ListEnvironments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEnvironmentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).ListEnvironments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/ListEnvironments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).ListEnvironments(ctx, req.(*ListEnvironmentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_GetEnvironment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEnvironmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).GetEnvironment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/GetEnvironment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).GetEnvironment(ctx, req.(*GetEnvironmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_DeleteEnvironment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEnvironmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).DeleteEnvironment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/DeleteEnvironment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).DeleteEnvironment(ctx, req.(*DeleteEnvironmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_CreateFeatureToggle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFeatureToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).CreateFeatureToggle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/CreateFeatureToggle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).CreateFeatureToggle(ctx, req.(*CreateFeatureToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_ListFeatureToggles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFeatureToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).ListFeatureToggles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/ListFeatureToggles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).ListFeatureToggles(ctx, req.(*ListFeatureToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_GetFeatureToggle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFeatureToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).GetFeatureToggle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/GetFeatureToggle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).GetFeatureToggle(ctx, req.(*GetFeatureToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_GetFeatureToggleHistoryForEnvironment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFeatureToggleHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).GetFeatureToggleHistoryForEnvironment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/GetFeatureToggleHistoryForEnvironment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).GetFeatureToggleHistoryForEnvironment(ctx, req.(*GetFeatureToggleHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_UpdateFeatureToggle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFeatureToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).UpdateFeatureToggle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/UpdateFeatureToggle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).UpdateFeatureToggle(ctx, req.(*UpdateFeatureToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dashboard_DeleteFeatureToggle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFeatureToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DashboardServer).DeleteFeatureToggle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.Dashboard/DeleteFeatureToggle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DashboardServer).DeleteFeatureToggle(ctx, req.(*DeleteFeatureToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Dashboard_ServiceDesc is the grpc.ServiceDesc for Dashboard service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Dashboard_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dashboard.Dashboard",
	HandlerType: (*DashboardServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _Dashboard_GetUser_Handler,
		},
		{
			MethodName: "CreateProject",
			Handler:    _Dashboard_CreateProject_Handler,
		},
		{
			MethodName: "ListProjects",
			Handler:    _Dashboard_ListProjects_Handler,
		},
		{
			MethodName: "GetProject",
			Handler:    _Dashboard_GetProject_Handler,
		},
		{
			MethodName: "DeleteProject",
			Handler:    _Dashboard_DeleteProject_Handler,
		},
		{
			MethodName: "CreateProjectInvite",
			Handler:    _Dashboard_CreateProjectInvite_Handler,
		},
		{
			MethodName: "ListProjectInvites",
			Handler:    _Dashboard_ListProjectInvites_Handler,
		},
		{
			MethodName: "ListUserInvites",
			Handler:    _Dashboard_ListUserInvites_Handler,
		},
		{
			MethodName: "ListProjectMembers",
			Handler:    _Dashboard_ListProjectMembers_Handler,
		},
		{
			MethodName: "GetProjectInvite",
			Handler:    _Dashboard_GetProjectInvite_Handler,
		},
		{
			MethodName: "UpdateProjectInvite",
			Handler:    _Dashboard_UpdateProjectInvite_Handler,
		},
		{
			MethodName: "CreateEnvironment",
			Handler:    _Dashboard_CreateEnvironment_Handler,
		},
		{
			MethodName: "ListEnvironments",
			Handler:    _Dashboard_ListEnvironments_Handler,
		},
		{
			MethodName: "GetEnvironment",
			Handler:    _Dashboard_GetEnvironment_Handler,
		},
		{
			MethodName: "DeleteEnvironment",
			Handler:    _Dashboard_DeleteEnvironment_Handler,
		},
		{
			MethodName: "CreateFeatureToggle",
			Handler:    _Dashboard_CreateFeatureToggle_Handler,
		},
		{
			MethodName: "ListFeatureToggles",
			Handler:    _Dashboard_ListFeatureToggles_Handler,
		},
		{
			MethodName: "GetFeatureToggle",
			Handler:    _Dashboard_GetFeatureToggle_Handler,
		},
		{
			MethodName: "GetFeatureToggleHistoryForEnvironment",
			Handler:    _Dashboard_GetFeatureToggleHistoryForEnvironment_Handler,
		},
		{
			MethodName: "UpdateFeatureToggle",
			Handler:    _Dashboard_UpdateFeatureToggle_Handler,
		},
		{
			MethodName: "DeleteFeatureToggle",
			Handler:    _Dashboard_DeleteFeatureToggle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dashboard.proto",
}
