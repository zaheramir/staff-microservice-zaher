package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"

	spb "github.com/BetterGR/staff-microservice/protos"
	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog"
)

// MockClaims overrides Claims behavior for testing.
type MockClaims struct {
	ms.Claims
}

// Always return true for HasRole.
func (m MockClaims) HasRole(_ string) bool {
	return true
}

// Always return "staff" for GetRole.
func (m MockClaims) GetRole() string {
	return "test-role"
}

// TestStaffServer wraps StaffServer for testing.
type TestStaffServer struct {
	*StaffServer
}

func TestMain(m *testing.M) {
	// Load .env file.
	cmd := exec.Command("cat", "../.env")

	output, err := cmd.Output()
	if err != nil {
		panic("Error reading .env file: " + err.Error())
	}

	// Set environment variables.
	for _, line := range strings.Split(string(output), "\n") {
		if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				// Remove quotes from the value if they exist.
				value := strings.Trim(parts[1], `"'`)
				os.Setenv(parts[0], value)
			}
		}
	}

	// Run tests and capture the result.
	result := m.Run()

	if result == 0 {
		klog.Info("\n\n [Summary] All tests passed.")
	} else {
		klog.Errorf("\n\n [Summary] Some tests failed. number of tests that failed: %d", result)
	}

	// Exit with the test result code.
	os.Exit(result)
}

func createTestStaffMember() *spb.StaffMember {
	return &spb.StaffMember{
		StaffID:     uuid.New().String(),
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@example.com",
		PhoneNumber: "1234567890",
	}
}

func startTestServer() (*grpc.Server, net.Listener, *TestStaffServer, error) {
	server, err := initStaffMicroserviceServer()
	if err != nil {
		return nil, nil, nil, err
	}

	server.Claims = MockClaims{}
	testServer := &TestStaffServer{StaffServer: server}
	grpcServer := grpc.NewServer()
	spb.RegisterStaffServiceServer(grpcServer, testServer)

	listener, err := net.Listen(connectionProtocol, "localhost:"+os.Getenv("GRPC_PORT"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to listen on port %s: %w", os.Getenv("GRPC_PORT"), err)
	}

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic("Failed to serve: " + err.Error())
		}
	}()

	return grpcServer, listener, testServer, nil
}

func setupClient(t *testing.T) spb.StaffServiceClient {
	t.Helper()

	grpcServer, listener, _, err := startTestServer()
	require.NoError(t, err)
	t.Cleanup(func() {
		grpcServer.Stop()
	})

	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() {
		conn.Close()
	})

	return spb.NewStaffServiceClient(conn)
}

func TestGetStaffMemberFound(t *testing.T) {
	client := setupClient(t)
	staffMember := createTestStaffMember()
	_, err := client.CreateStaffMember(t.Context(),
		&spb.CreateStaffMemberRequest{StaffMember: staffMember, Token: "test-token"})
	require.NoError(t, err)

	req := &spb.GetStaffMemberRequest{StaffID: staffMember.GetStaffID(), Token: "test-token"}
	resp, err := client.GetStaffMember(t.Context(), req)
	require.NoError(t, err)
	assert.Equal(t, staffMember.GetStaffID(), resp.GetStaffMember().GetStaffID())

	// Cleanup.
	_, _ = client.DeleteStaffMember(t.Context(),
		&spb.DeleteStaffMemberRequest{StaffID: staffMember.GetStaffID(), Token: "test-token"})
}

func TestGetStaffMemberNotFound(t *testing.T) {
	client := setupClient(t)
	req := &spb.GetStaffMemberRequest{StaffID: "non-existent-id", Token: "test-token"}

	_, err := client.GetStaffMember(t.Context(), req)
	assert.Error(t, err)
}

func TestCreateStaffMemberSuccessful(t *testing.T) {
	client := setupClient(t)
	staffMember := createTestStaffMember()
	req := &spb.CreateStaffMemberRequest{StaffMember: staffMember, Token: "test-token"}

	resp, err := client.CreateStaffMember(t.Context(), req)
	require.NoError(t, err)
	assert.Equal(t, resp.GetStaffMember().GetEmail(), staffMember.GetEmail())

	// Cleanup.
	_, _ = client.DeleteStaffMember(t.Context(),
		&spb.DeleteStaffMemberRequest{StaffID: staffMember.GetStaffID(), Token: "test-token"})
}

func TestCreateStaffMemberFailureOnDuplicate(t *testing.T) {
	client := setupClient(t)
	staffMember := createTestStaffMember()
	_, err := client.CreateStaffMember(t.Context(),
		&spb.CreateStaffMemberRequest{StaffMember: staffMember, Token: "test-token"})
	require.NoError(t, err)

	req := &spb.CreateStaffMemberRequest{StaffMember: staffMember, Token: "test-token"}
	_, err = client.CreateStaffMember(t.Context(), req)
	require.Error(t, err)

	// Cleanup.
	_, _ = client.DeleteStaffMember(t.Context(),
		&spb.DeleteStaffMemberRequest{StaffID: staffMember.GetStaffID(), Token: "test-token"})
}

func TestUpdateStaffMemberSuccessful(t *testing.T) {
	client := setupClient(t)
	staffMember := createTestStaffMember()
	_, err := client.CreateStaffMember(t.Context(),
		&spb.CreateStaffMemberRequest{StaffMember: staffMember, Token: "test-token"})
	require.NoError(t, err)

	// Modify staff member.
	staffMember.FirstName = "UpdatedName"
	req := &spb.UpdateStaffMemberRequest{StaffMember: staffMember, Token: "test-token"}

	resp, err := client.UpdateStaffMember(t.Context(), req)
	require.NoError(t, err)
	assert.Equal(t, resp.GetStaffMember().GetFirstName(), staffMember.GetFirstName())
	// Cleanup.

	_, _ = client.DeleteStaffMember(t.Context(),
		&spb.DeleteStaffMemberRequest{StaffID: staffMember.GetStaffID(), Token: "test-token"})
}

func TestUpdateStaffMemberFailureForNonExistentStaffMember(t *testing.T) {
	client := setupClient(t)
	staffMember := createTestStaffMember()
	staffMember.StaffID = "non-existent-id"
	req := &spb.UpdateStaffMemberRequest{StaffMember: staffMember, Token: "test-token"}

	_, err := client.UpdateStaffMember(t.Context(), req)
	assert.Error(t, err)
}

func TestDeleteStaffMemberSuccessful(t *testing.T) {
	client := setupClient(t)
	staffMember := createTestStaffMember()
	_, err := client.CreateStaffMember(t.Context(),
		&spb.CreateStaffMemberRequest{StaffMember: staffMember, Token: "test-token"})
	require.NoError(t, err)

	req := &spb.DeleteStaffMemberRequest{StaffID: staffMember.GetStaffID(), Token: "test-token"}
	_, err = client.DeleteStaffMember(t.Context(), req)
	assert.NoError(t, err)
}

func TestDeleteStaffMemberFailureForNonExistentStaffMember(t *testing.T) {
	client := setupClient(t)
	req := &spb.DeleteStaffMemberRequest{StaffID: "non-existent-id", Token: "test-token"}

	_, err := client.DeleteStaffMember(t.Context(), req)
	assert.Error(t, err)
}
