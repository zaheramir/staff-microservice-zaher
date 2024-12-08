package main

import (
	"context"
	"log"
	"net"

	pb "github.com/BetterGR/staff-microservice/staff_protobuf"
	"google.golang.org/grpc"
)

// server is used to implement pb.StaffServiceServer
type staffServer struct {
	pb.UnimplementedStaffServiceServer
}

const (
	port = ":50051"
)

/*// GetStudentGrade method
func (s *gradesServer) GetStudentGrade(ctx context.Context, req *gpb.GradeRequest) (*gpb.GradeReply, error) {
	log.Printf("Recevied", req.GetStudentId())
	return &gpb.GradeReply{Grade: "100", Course: "test"}, nil
}
*/

// GetStaffMember retrieves a specific staff member
func (s *staffServer) GetStaffMember(ctx context.Context, req *pb.GetStaffMemberRequest) (
	*pb.GetStaffMemberResponse, error) {

	log.Printf("Received GetStaffMember request: %v", req)
	//TODO: implement the method
	return &pb.GetStaffMemberResponse{}, nil
}

// GetCoursesList retrieves all courses assigned to a staff member
func (s *staffServer) GetCoursesList(ctx context.Context, req *pb.GetCoursesListRequest) (
	*pb.GetCoursesListResponse, error) {
	log.Printf("Received GetCoursesList request: %v", req)
	//TODO: implement the method
	return &pb.GetCoursesListResponse{}, nil
}

// CreateStaffMember creates a new staff member
func (s *staffServer) CreateStaffMember(ctx context.Context, req *pb.CreateStaffMemberRequest) (
	*pb.CreateStaffMemberResponse, error) {
	log.Printf("Received CreateStaffMember request: %v", req)
	//TODO: implement the method
	return &pb.CreateStaffMemberResponse{}, nil
}

// UpdateStaffMember updates details of an existing staff member
func (s *staffServer) UpdateStaffMember(ctx context.Context, req *pb.UpdateStaffMemberRequest) (
	*pb.UpdateStaffMemberResponse, error) {
	log.Printf("Received UpdateStaffMember request: %v", req)
	//TODO: implement the method
	return &pb.UpdateStaffMemberResponse{}, nil
}

// DeleteStaffMember deletes a specific staff member
func (s *staffServer) DeleteStaffMember(ctx context.Context, req *pb.DeleteStaffMemberRequest) (
	*pb.DeleteStaffMemberResponse, error) {
	log.Printf("Received DeleteStaffMember request: %v", req)
	//TODO: implement the method
	return &pb.DeleteStaffMemberResponse{}, nil
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStaffServiceServer(grpcServer, &staffServer{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
