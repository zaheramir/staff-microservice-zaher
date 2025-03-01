// main package to be able to run the StaffServer for now
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	spb "github.com/BetterGR/staff-microservice/protos"
	ms "github.com/TekClinic/MicroService-Lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	// define address.
	connectionProtocol = "tcp"
	// Debugging logs.
	logLevelDebug = 5
)

// StaffServer is an implementation of GRPC Staff microservice.
type StaffServer struct {
	ms.BaseServiceServer
	db *Database
	spb.UnimplementedStaffServiceServer
	Claims ms.Claims
}

// VerifyToken returns the injected Claims instead of the default.
func (s *StaffServer) VerifyToken(ctx context.Context, token string) error {
	if s.Claims != nil {
		return nil
	}

	// Default behavior.
	if _, err := s.BaseServiceServer.VerifyToken(ctx, token); err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}

	return nil
}

func initStaffMicroserviceServer() (*StaffServer, error) {
	base, err := ms.CreateBaseServiceServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create base service: %w", err)
	}

	database, err := InitializeDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &StaffServer{
		BaseServiceServer:               base,
		db:                              database,
		UnimplementedStaffServiceServer: spb.UnimplementedStaffServiceServer{},
	}, nil
}

// GetStaffMember search for the StaffMember that corresponds to the given id and returns them.
func (s *StaffServer) GetStaffMember(ctx context.Context,
	req *spb.GetStaffMemberRequest,
) (*spb.GetStaffMemberResponse, error) {
	if err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received GetStaffMember request", "staffId", req.GetStaffID())

	staff, err := s.db.GetStaffMember(ctx, req.GetStaffID())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "staff member not found: %v", err)
	}

	staffMember := &spb.StaffMember{
		StaffID:     staff.StaffID,
		FirstName:   staff.FirstName,
		LastName:    staff.LastName,
		Email:       staff.Email,
		PhoneNumber: staff.PhoneNumber,
		Title:       staff.Title,
		Office:      staff.Office,
	}

	return &spb.GetStaffMemberResponse{StaffMember: staffMember}, nil
}

// CreateStaffMember creates a new StaffMember with the given details and returns them.
func (s *StaffServer) CreateStaffMember(ctx context.Context,
	req *spb.CreateStaffMemberRequest,
) (*spb.CreateStaffMemberResponse, error) {
	if err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received CreateStaffMember request",
		"firstName", req.GetStaffMember().GetFirstName(), "secondName", req.GetStaffMember().GetLastName())

	if _, err := s.db.AddStaffMember(ctx, req.GetStaffMember()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create staff member: %v", err)
	}

	return &spb.CreateStaffMemberResponse{StaffMember: req.GetStaffMember()}, nil
}

// UpdateStaffMember updates the given StaffMember and returns them after the update.
func (s *StaffServer) UpdateStaffMember(ctx context.Context,
	req *spb.UpdateStaffMemberRequest,
) (*spb.UpdateStaffMemberResponse, error) {
	if err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received UpdateStaffMember request",
		"firstName", req.GetStaffMember().GetFirstName(), "secondName", req.GetStaffMember().GetLastName())

	updatedStaff, err := s.db.UpdateStaffMember(ctx, req.GetStaffMember())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update staff member: %v", err)
	}

	staff := &spb.StaffMember{
		StaffID:     updatedStaff.StaffID,
		FirstName:   updatedStaff.FirstName,
		LastName:    updatedStaff.LastName,
		Email:       updatedStaff.Email,
		PhoneNumber: updatedStaff.PhoneNumber,
		Title:       updatedStaff.Title,
		Office:      updatedStaff.Office,
	}

	return &spb.UpdateStaffMemberResponse{StaffMember: staff}, nil
}

// DeleteStaffMember deletes the StaffMember from the system.
func (s *StaffServer) DeleteStaffMember(ctx context.Context,
	req *spb.DeleteStaffMemberRequest,
) (*spb.DeleteStaffMemberResponse, error) {
	if err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(logLevelDebug).Info("Received DeleteStaffMember request", "staffId", req.GetStaffID())

	if err := s.db.DeleteStaffMember(ctx, req.GetStaffID()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete staff member: %v", err)
	}

	logger.V(logLevelDebug).Info("Deleted", "staffId", req.GetStaffID())

	return &spb.DeleteStaffMemberResponse{}, nil
}

// main StaffServer function.
func main() {
	// init klog
	klog.InitFlags(nil)
	flag.Parse()

	// init the StaffServer
	server, err := initStaffMicroserviceServer()
	if err != nil {
		klog.Fatalf("Failed to init StaffServer: %v", err)
	}

	// create a listener on port 'address'
	address := os.Getenv("STAFF_PORT")

	lis, err := net.Listen(connectionProtocol, address)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	klog.V(logLevelDebug).Info("Starting StaffServer on port: ", address)
	// create a grpc StaffServer
	grpcServer := grpc.NewServer()
	spb.RegisterStaffServiceServer(grpcServer, server)

	// serve the grpc StaffServer
	if err := grpcServer.Serve(lis); err != nil {
		klog.Fatalf("Failed to serve: %v", err)
	}
}
