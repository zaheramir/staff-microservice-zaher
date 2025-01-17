package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	staffProtos "github.com/BetterGR/staff-microservice/protos"
	ms "github.com/TekClinic/MicroService-Lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	connectionProtocol = "tcp"
	traceVerbosity     = 5
)

// server is used to implement staffProtos.StaffServiceServer.
type staffServer struct {
	ms.BaseServiceServer
	staffProtos.UnimplementedStaffServiceServer
}

// createStaffMicroserviceServer initiates the mslib for staff microservice.
func createStaffMicroserviceServer() (*staffServer, error) {
	base, err := ms.CreateBaseServiceServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create base service: %w", err)
	}

	return &staffServer{
		BaseServiceServer:               base,
		UnimplementedStaffServiceServer: staffProtos.UnimplementedStaffServiceServer{},
	}, nil
}

// GetStaffMember retrieves a specific staff member.
func (s *staffServer) GetStaffMember(ctx context.Context, req *staffProtos.GetStaffMemberRequest) (
	*staffProtos.GetStaffMemberResponse, error,
) {
	_, err := s.VerifyToken(ctx, req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	staffMemberID := req.GetId()
	logger.V(traceVerbosity).Info("Received GetStaffMember request", "staffMemberID", staffMemberID)
	// TODO: implement the method
	// for the sake of mslib integration, we throw an error as a temporary return, because method is not implemented yet
	return nil, status.Errorf(codes.Unimplemented, "method GetStaffMember not implemented")
}

// GetCoursesList retrieves all courses assigned to a staff member.
func (s *staffServer) GetCoursesList(ctx context.Context, req *staffProtos.GetCoursesListRequest) (
	*staffProtos.GetCoursesListResponse, error,
) {
	logger := klog.FromContext(ctx)
	firstName := req.GetStaffMember().GetFirstName()
	secondName := req.GetStaffMember().GetSecondName()
	semester := req.GetSemester()
	logger.V(traceVerbosity).Info("Received GetCoursesList request",
		"firstName", firstName, "secondName", secondName, "semester", semester)
	// TODO: implement the method
	// for the sake of mslib integration, we throw an error as a temporary return, because method is not implemented yet
	return nil, status.Errorf(codes.Unimplemented, "method GetCoursesList not implemented")
}

// CreateStaffMember creates a new staff member.
func (s *staffServer) CreateStaffMember(ctx context.Context, req *staffProtos.CreateStaffMemberRequest) (
	*staffProtos.CreateStaffMemberResponse, error,
) {
	logger := klog.FromContext(ctx)
	firstName := req.GetStaffMember().GetFirstName()
	secondName := req.GetStaffMember().GetSecondName()
	logger.V(traceVerbosity).Info("Received CreateStaffMember request",
		"firstName", firstName, "secondName", secondName)
	// TODO: implement the method
	// for the sake of mslib integration, we throw an error as a temporary return, because method is not implemented yet
	return nil, status.Errorf(codes.Unimplemented, "method CreateStaffMember not implemented")
}

// UpdateStaffMember updates details of an existing staff member.
func (s *staffServer) UpdateStaffMember(ctx context.Context, req *staffProtos.UpdateStaffMemberRequest) (
	*staffProtos.UpdateStaffMemberResponse, error,
) {
	logger := klog.FromContext(ctx)
	firstName := req.GetStaffMember().GetFirstName()
	secondName := req.GetStaffMember().GetSecondName()
	logger.V(traceVerbosity).Info("Received UpdateStaffMember request",
		"firstName", firstName, "secondName", secondName)

	// TODO: implement the method
	// for the sake of mslib integration, we throw an error as a temporary return, because method is not implemented yet
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStaffMember not implemented")
}

// DeleteStaffMember deletes a specific staff member.
func (s *staffServer) DeleteStaffMember(ctx context.Context, req *staffProtos.DeleteStaffMemberRequest) (
	*staffProtos.DeleteStaffMemberResponse, error,
) {
	logger := klog.FromContext(ctx)
	staffMemberID := req.GetStaffMember().GetId()
	logger.V(traceVerbosity).Info("Received DeleteStaffMember request", "staffMemberID", staffMemberID)
	// TODO: implement the method
	// for the sake of mslib integration, we throw an error as a temporary return, because method is not implemented yet
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStaffMember not implemented")
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// init the StaffServer
	server, err := createStaffMicroserviceServer()
	if err != nil {
		klog.Error("Failed to init StaffServer", err)
	}

	// create a listener on the port specified in file .env
	listener, err := net.Listen(connectionProtocol, "localhost:"+os.Getenv("STAFF_PORT"))
	if err != nil {
		klog.Error("Failed to listen:", err)
	}

	// create a grpc StaffServer
	grpcServer := grpc.NewServer()
	staffProtos.RegisterStaffServiceServer(grpcServer, server)

	if err := grpcServer.Serve(listener); err != nil {
		klog.Error("Failed to serve:", err)
	}
}
