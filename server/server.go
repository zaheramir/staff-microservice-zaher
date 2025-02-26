package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	spb "github.com/BetterGR/staff-microservice/protos"
	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

const (
	connectionProtocol = "tcp"
	traceVerbosity     = 0
)

// staffServer implements staffProtos.StaffServiceServer.
type staffServer struct {
	ms.BaseServiceServer
	spb.UnimplementedStaffServiceServer
	db *Database
}

// createStaffMicroserviceServer initiates the mslib for staff microservice.
func createStaffMicroserviceServer() (*staffServer, error) {
	base, err := ms.CreateBaseServiceServer()
	if err != nil {
		return nil, fmt.Errorf("failed to create base service: %w", err)
	}

	database, err := InitializeDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &staffServer{
		BaseServiceServer:               base,
		UnimplementedStaffServiceServer: spb.UnimplementedStaffServiceServer{},
		db:                              database,
	}, nil
}

// GetStaffMember retrieves a specific staff member.
func (s *staffServer) GetStaffMember(ctx context.Context, req *spb.GetStaffMemberRequest) (
	*spb.GetStaffMemberResponse, error,
) {
	if _, err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(traceVerbosity).Info("Received GetStaffMember request", "staffMemberID", req.GetStaffID())

	staffMember, err := s.db.GetStaff(ctx, req.GetStaffID())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "staff member not found: %v", err)
	}

	return &spb.GetStaffMemberResponse{StaffMember: staffMember}, nil
}

// CreateStaffMember creates a new staff member.
func (s *staffServer) CreateStaffMember(ctx context.Context, req *spb.CreateStaffMemberRequest) (
	*spb.CreateStaffMemberResponse, error,
) {
	if _, err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(traceVerbosity).Info("Received CreateStaffMember request", "firstName", req.GetStaffMember().GetFirstName())

	if err := s.db.AddStaff(ctx, req.GetStaffMember()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create staff member: %v", err)
	}

	return &spb.CreateStaffMemberResponse{StaffMember: req.GetStaffMember()}, nil
}

// UpdateStaffMember updates details of an existing staff member.
func (s *staffServer) UpdateStaffMember(ctx context.Context, req *spb.UpdateStaffMemberRequest) (
	*spb.UpdateStaffMemberResponse, error,
) {
	if _, err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	logger.V(traceVerbosity).Info("Received UpdateStaffMember request", "staffMemberID", req.GetStaffMember().GetStaffID())

	if err := s.db.UpdateStaff(ctx, req.GetStaffMember()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update staff member: %v", err)
	}

	return &spb.UpdateStaffMemberResponse{StaffMember: req.GetStaffMember()}, nil
}

// DeleteStaffMember deletes a specific staff member.
func (s *staffServer) DeleteStaffMember(ctx context.Context, req *spb.DeleteStaffMemberRequest) (
	*spb.DeleteStaffMemberResponse, error,
) {
	if _, err := s.VerifyToken(ctx, req.GetToken()); err != nil {
		return nil, fmt.Errorf("authentication failed: %w",
			status.Error(codes.Unauthenticated, err.Error()))
	}

	logger := klog.FromContext(ctx)
	id := req.GetStaffID()
	logger.V(traceVerbosity).Info("Received DeleteStaffMember request", "staffMemberID", id)

	if err := s.db.DeleteStaff(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete staff member: %v", err)
	}

	return &spb.DeleteStaffMemberResponse{}, nil
}

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		klog.Fatalf("Error loading .env file")
	}

	// init the StaffServer
	server, err := createStaffMicroserviceServer()
	if err != nil {
		klog.Error("Failed to init StaffServer", err)
	}

	// create a listener on the port specified in file .env
	address := os.Getenv("GRPC_PORT")

	listener, err := net.Listen(connectionProtocol, address)
	if err != nil {
		klog.Error("Failed to listen:", err)
	}

	klog.Info("Starting StudentsServer on port: ", address)
	// create a grpc StaffServer
	grpcServer := grpc.NewServer()
	spb.RegisterStaffServiceServer(grpcServer, server)

	if err := grpcServer.Serve(listener); err != nil {
		klog.Error("Failed to serve:", err)
	}
}
