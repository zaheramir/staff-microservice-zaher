package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	spb "github.com/BetterGR/staff-microservice/protos"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"k8s.io/klog/v2"
)

// Database represents the PostgreSQL database connection.
type Database struct {
	db *bun.DB
}

var (
	ErrStaffMemberNil      = errors.New("staff member is nil")
	ErrStaffMemberIDEmpty  = errors.New("staff member ID is empty")
	ErrStaffMemberNotFound = errors.New("staff member not found")
)

// InitializeDatabase ensures that the database exists and initializes the schema.
func InitializeDatabase() (*Database, error) {
	createDatabaseIfNotExists()

	database, err := ConnectDB()
	if err != nil {
		return nil, err
	}

	if err := database.createSchemaIfNotExists(context.Background()); err != nil {
		klog.Fatalf("Failed to create schema: %v", err)
	}

	return database, nil
}

// createDatabaseIfNotExists ensures the database exists.
func createDatabaseIfNotExists() {
	dsn := os.Getenv("DSN")
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))

	sqldb := sql.OpenDB(connector)
	defer sqldb.Close()

	ctx := context.Background()
	dbName := os.Getenv("DP_NAME")
	query := "SELECT 1 FROM pg_database WHERE datname = $1;"

	var exists int

	err := sqldb.QueryRowContext(ctx, query, dbName).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		klog.Fatalf("Failed to check db existence: %v", err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		if _, err = sqldb.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s;", dbName)); err != nil {
			klog.Fatalf("Failed to create database: %v", err)
		}

		klog.V(logLevelDebug).Infof("Database %s created successfully.", dbName)
	} else {
		klog.V(logLevelDebug).Infof("Database %s already exists.", dbName)
	}
}

// ConnectDB connects to the database.
func ConnectDB() (*Database, error) {
	dsn := os.Getenv("DSN")
	connector := pgdriver.NewConnector(pgdriver.WithDSN(dsn))
	sqldb := sql.OpenDB(connector)
	database := bun.NewDB(sqldb, pgdialect.New())

	// Test the connection.
	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	klog.V(logLevelDebug).Info("Connected to PostgreSQL database.")

	return &Database{db: database}, nil
}

// createSchemaIfNotExists creates the database schema if it doesn't exist.
func (d *Database) createSchemaIfNotExists(ctx context.Context) error {
	models := []interface{}{
		(*StaffMember)(nil),
	}

	for _, model := range models {
		if _, err := d.db.NewCreateTable().IfNotExists().Model(model).Exec(ctx); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	klog.V(logLevelDebug).Info("Database schema initialized.")

	return nil
}

// StaffMember represents the staff_member table.
type StaffMember struct {
	StaffID     string    `bun:"staff_id,unique,pk,notnull"`
	FirstName   string    `bun:"first_name,notnull"`
	LastName    string    `bun:"last_name,notnull"`
	Email       string    `bun:"email,unique,notnull"`
	PhoneNumber string    `bun:"phone_number,unique,notnull"`
	Title       string    `bun:"title"`
	Office      string    `bun:"office"`
	CreatedAt   time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,default:current_timestamp"`
}

// AddStaffMember adds a new staff member.
func (d *Database) AddStaffMember(ctx context.Context, staff *spb.StaffMember) (*StaffMember, error) {
	if staff == nil {
		return nil, fmt.Errorf("%w", ErrStaffMemberNil)
	}

	newStaffMember := &StaffMember{
		StaffID:     staff.GetStaffID(),
		FirstName:   staff.GetFirstName(),
		LastName:    staff.GetLastName(),
		Email:       staff.GetEmail(),
		PhoneNumber: staff.GetPhoneNumber(),
		Title:       staff.GetTitle(),
		Office:      staff.GetOffice(),
	}

	if _, err := d.db.NewInsert().Model(newStaffMember).Exec(ctx); err != nil {
		return nil, fmt.Errorf("failed to add staff member: %w", err)
	}

	return newStaffMember, nil
}

// GetStaffMember retrieves a staff member by ID.
func (d *Database) GetStaffMember(ctx context.Context, staffID string) (*StaffMember, error) {
	if staffID == "" {
		return nil, fmt.Errorf("%w", ErrStaffMemberIDEmpty)
	}

	staffMember := new(StaffMember)
	if err := d.db.NewSelect().Model(staffMember).Where("staff_id = ?", staffID).Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to get staff member: %w", err)
	}

	return staffMember, nil
}

// UpdateStaffMember updates an existing staff member.
func (d *Database) UpdateStaffMember(ctx context.Context, staff *spb.StaffMember) (*StaffMember, error) {
	if staff == nil {
		return nil, fmt.Errorf("%w", ErrStaffMemberNil)
	}

	if staff.GetStaffID() == "" {
		return nil, fmt.Errorf("%w", ErrStaffMemberIDEmpty)
	}

	// get the existing staff member
	existingStaffMember := &StaffMember{StaffID: staff.GetStaffID()}
	if err := d.db.NewSelect().Model(existingStaffMember).WherePK().Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to get staff member: %w", err)
	}

	// Update the fields.
	updateField := func(field *string, newValue string) {
		if newValue != "" {
			*field = newValue
		}
	}

	updateField(&existingStaffMember.FirstName, staff.GetFirstName())
	updateField(&existingStaffMember.LastName, staff.GetLastName())
	updateField(&existingStaffMember.Email, staff.GetEmail())
	updateField(&existingStaffMember.PhoneNumber, staff.GetPhoneNumber())
	updateField(&existingStaffMember.Title, staff.GetTitle())
	updateField(&existingStaffMember.Office, staff.GetOffice())

	if _, err := d.db.NewUpdate().Model(existingStaffMember).WherePK().Exec(ctx); err != nil {
		return nil, fmt.Errorf("failed to update staff member: %w", err)
	}

	return existingStaffMember, nil
}

// DeleteStaffMember deletes a staff member by ID.
func (d *Database) DeleteStaffMember(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w", ErrStaffMemberIDEmpty)
	}

	res, err := d.db.NewDelete().Model((*StaffMember)(nil)).Where("staff_id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete staff member: %w", err)
	}

	if num, _ := res.RowsAffected(); num == 0 {
		return fmt.Errorf("%w", ErrStaffMemberNotFound)
	}

	return nil
}
