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

		klog.Infof("Database %s created successfully.", dbName)
	} else {
		klog.Infof("Database %s already exists.", dbName)
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

	klog.Info("Connected to PostgreSQL database.")

	return &Database{db: database}, nil
}

// createSchemaIfNotExists creates the database schema if it doesn't exist.
func (d *Database) createSchemaIfNotExists(ctx context.Context) error {
	models := []interface{}{
		(*Staff)(nil),
	}

	for _, model := range models {
		if _, err := d.db.NewCreateTable().IfNotExists().Model(model).Exec(ctx); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	klog.Info("Database schema initialized.")

	return nil
}

// Staff represents the staff table.
type Staff struct {
	StaffID     string    `bun:"staff_id,unique,notnull"`
	FirstName   string    `bun:"first_name,notnull"`
	LastName    string    `bun:"last_name,notnull"`
	Email       string    `bun:"email,unique,notnull"`
	PhoneNumber string    `bun:"phone_number,unique,notnull"`
	Title       string    `bun:"title,notnull"`
	Office      string    `bun:"office,notnull"`
	CreatedAt   time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"updated_at,default:current_timestamp"`
}

// AddStaff adds a new staff member.
func (d *Database) AddStaff(ctx context.Context, staff *spb.StaffMember) error {
	_, err := d.db.NewInsert().Model(&Staff{
		StaffID:     staff.GetStaffID(),
		FirstName:   staff.GetFirstName(),
		LastName:    staff.GetSecondName(),
		Email:       staff.GetEmail(),
		PhoneNumber: staff.GetPhoneNumber(),
		Title:       staff.GetTitle(),
		Office:      staff.GetOffice(),
	}).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add staff: %w", err)
	}

	return nil
}

// GetStaff retrieves a staff member by ID.
func (d *Database) GetStaff(ctx context.Context, id string) (*spb.StaffMember, error) {
	staffMember := new(Staff)
	if err := d.db.NewSelect().Model(staffMember).Where("staff_id = ?", id).Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	return &spb.StaffMember{
		StaffID:     staffMember.StaffID,
		FirstName:   staffMember.FirstName,
		SecondName:  staffMember.LastName,
		Email:       staffMember.Email,
		PhoneNumber: staffMember.PhoneNumber,
		Title:       staffMember.Title,
		Office:      staffMember.Office,
	}, nil
}

// UpdateStaff updates an existing staff member.
func (d *Database) UpdateStaff(ctx context.Context, staff *spb.StaffMember) error {
	_, err := d.db.NewUpdate().Model(&Staff{
		StaffID:     staff.GetStaffID(),
		FirstName:   staff.GetFirstName(),
		LastName:    staff.GetSecondName(),
		Email:       staff.GetEmail(),
		PhoneNumber: staff.GetPhoneNumber(),
		Title:       staff.GetTitle(),
		Office:      staff.GetOffice(),
	}).Where("staff_id = ?", staff.GetStaffID()).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update staff: %w", err)
	}

	return nil
}

// DeleteStaff deletes a staff member by ID.
func (d *Database) DeleteStaff(ctx context.Context, id string) error {
	_, err := d.db.NewDelete().Model((*Staff)(nil)).Where("staff_id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete staff: %w", err)
	}

	return nil
}
