package main

import (
	"context"
	pb "github.com/BetterGR/staff-microservice/staff_protobuf"
	"google.golang.org/grpc"
	"k8s.io/klog/v2"
	"net"
)

const (
	tcpProtocol = "tcp"
	port        = "50051"
)

// server is used to implement pb.StaffServiceServer.
type staffServer struct {
	pb.UnimplementedStaffServiceServer
}

// GetStaffMember retrieves a specific staff member.
func (s *staffServer) GetStaffMember(ctx context.Context, req *pb.GetStaffMemberRequest) (
	*pb.GetStaffMemberResponse, error) {

	staffMemberId := req.GetId()
	logger := klog.FromContext(ctx)
	logger.Info("Received GetStaffMember request", "staffMemberId", staffMemberId)
	//TODO: implement the method
	return &pb.GetStaffMemberResponse{}, nil
}

// GetCoursesList retrieves all courses assigned to a staff member.
func (s *staffServer) GetCoursesList(ctx context.Context, req *pb.GetCoursesListRequest) (
	*pb.GetCoursesListResponse, error) {

	firstName := req.GetStaffMember().GetFirstName()
	secondName := req.GetStaffMember().GetSecondName()
	semester := req.GetSemester()
	logger := klog.FromContext(ctx)
	logger.Info("Received GetCoursesList request",
		"firstName", firstName, "secondName", secondName, "semester", semester)
	//TODO: implement the method
	return &pb.GetCoursesListResponse{}, nil
}

// CreateStaffMember creates a new staff member.
func (s *staffServer) CreateStaffMember(ctx context.Context, req *pb.CreateStaffMemberRequest) (
	*pb.CreateStaffMemberResponse, error) {

	firstName := req.GetStaffMember().GetFirstName()
	secondName := req.GetStaffMember().GetSecondName()
	logger := klog.FromContext(ctx)
	logger.Info("Received CreateStaffMember request",
		"firstName", firstName, "secondName", secondName)
	//TODO: implement the method
	return &pb.CreateStaffMemberResponse{}, nil
}

// UpdateStaffMember updates details of an existing staff member.
func (s *staffServer) UpdateStaffMember(ctx context.Context, req *pb.UpdateStaffMemberRequest) (
	*pb.UpdateStaffMemberResponse, error) {

	firstName := req.GetStaffMember().GetFirstName()
	secondName := req.GetStaffMember().GetSecondName()
	logger := klog.FromContext(ctx)
	logger.Info("Received UpdateStaffMember request",
		"firstName", firstName, "secondName", secondName)

	//TODO: implement the method
	return &pb.UpdateStaffMemberResponse{}, nil
}

// DeleteStaffMember deletes a specific staff member.
func (s *staffServer) DeleteStaffMember(ctx context.Context, req *pb.DeleteStaffMemberRequest) (
	*pb.DeleteStaffMemberResponse, error) {

	staffMemberId := req.GetStaffMember().GetId()
	logger := klog.FromContext(ctx)
	logger.Info("Received DeleteStaffMember request", "staffMemberId", staffMemberId)
	//TODO: implement the method
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
