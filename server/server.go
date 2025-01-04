package main

import (
	"context"
	"net"

	pb "github.com/BetterGR/staff-microservice/protos"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

const (
	tcpProtocol    = "tcp"
	port           = "50051"
	traceVerbosity = 5
)

// server is used to implement pb.StaffServiceServer.
type staffServer struct {
	pb.UnimplementedStaffServiceServer
}

// GetStaffMember retrieves a specific staff member.
func (s *staffServer) GetStaffMember(ctx context.Context, req *pb.GetStaffMemberRequest) (
	*pb.GetStaffMemberResponse, error,
) {
	logger := klog.FromContext(ctx)
	staffMemberID := req.GetId()
	logger.V(traceVerbosity).Info("Received GetStaffMember request", "staffMemberID", staffMemberID)
	// TODO: implement the method
	return &pb.GetStaffMemberResponse{}, nil
}

// GetCoursesList retrieves all courses assigned to a staff member.
func (s *staffServer) GetCoursesList(ctx context.Context, req *pb.GetCoursesListRequest) (
	*pb.GetCoursesListResponse, error,
) {
	logger := klog.FromContext(ctx)
	firstName := req.GetStaffMember().GetFirstName()
	secondName := req.GetStaffMember().GetSecondName()
	semester := req.GetSemester()
	logger.V(traceVerbosity).Info("Received GetCoursesList request",
		"firstName", firstName, "secondName", secondName, "semester", semester)
	// TODO: implement the method
	return &pb.GetCoursesListResponse{}, nil
}

// CreateStaffMember creates a new staff member.
func (s *staffServer) CreateStaffMember(ctx context.Context, req *pb.CreateStaffMemberRequest) (
	*pb.CreateStaffMemberResponse, error,
) {
	logger := klog.FromContext(ctx)
	firstName := req.GetStaffMember().GetFirstName()
	secondName := req.GetStaffMember().GetSecondName()
	logger.V(traceVerbosity).Info("Received CreateStaffMember request",
		"firstName", firstName, "secondName", secondName)
	// TODO: implement the method
	return &pb.CreateStaffMemberResponse{}, nil
}

// UpdateStaffMember updates details of an existing staff member.
func (s *staffServer) UpdateStaffMember(ctx context.Context, req *pb.UpdateStaffMemberRequest) (
	*pb.UpdateStaffMemberResponse, error,
) {
	logger := klog.FromContext(ctx)
	firstName := req.GetStaffMember().GetFirstName()
	secondName := req.GetStaffMember().GetSecondName()
	logger.V(traceVerbosity).Info("Received UpdateStaffMember request",
		"firstName", firstName, "secondName", secondName)

	// TODO: implement the method
	return &pb.UpdateStaffMemberResponse{}, nil
}

// DeleteStaffMember deletes a specific staff member.
func (s *staffServer) DeleteStaffMember(ctx context.Context, req *pb.DeleteStaffMemberRequest) (
	*pb.DeleteStaffMemberResponse, error,
) {
	logger := klog.FromContext(ctx)
	staffMemberID := req.GetStaffMember().GetId()
	logger.V(traceVerbosity).Info("Received DeleteStaffMember request", "staffMemberID", staffMemberID)
	// TODO: implement the method
	return &pb.DeleteStaffMemberResponse{}, nil
}

func main() {
	listener, err := net.Listen(tcpProtocol, port)
	if err != nil {
		klog.Error("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStaffServiceServer(grpcServer, &staffServer{})

	if err := grpcServer.Serve(listener); err != nil {
		klog.Error("Failed to serve:", err)
	}
}
