package repository

import (
	"database/sql"
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// VolunteerRepository defines the interface for volunteer data operations
type VolunteerRepository interface {
	Create(volunteer *entity.Volunteer) error
	GetAll(page, limit int, status string) ([]entity.Volunteer, int, error)
	GetByID(id string) (*entity.Volunteer, error)
	GetByUserID(userID string) (*entity.Volunteer, error)
	HasPendingApplication(userID string) (bool, error)
	UpdateStatus(id string, status entity.ApplicationStatus) error
}

// volunteerRepository is the concrete implementation of VolunteerRepository
type volunteerRepository struct {
	db *sqlx.DB
}

// NewVolunteerRepository creates a new instance of VolunteerRepository
func NewVolunteerRepository(db *sqlx.DB) VolunteerRepository {
	return &volunteerRepository{db: db}
}

// Create inserts a new volunteer application into the database
func (r *volunteerRepository) Create(volunteer *entity.Volunteer) error {
	volunteer.CreatedAt = time.Now()
	volunteer.UpdatedAt = time.Now()

	query := `INSERT INTO volunteers (id, user_id, full_name, email, phone, date_of_birth, gender, city, occupation, interests, experience, status, created_at, updated_at) 
			  VALUES (:id, :user_id, :full_name, :email, :phone, :date_of_birth, :gender, :city, :occupation, :interests, :experience, :status, :created_at, :updated_at)`

	_, err := r.db.NamedExec(query, volunteer)
	return err
}

// GetAll retrieves all volunteer applications with pagination and optional status filtering
func (r *volunteerRepository) GetAll(page, limit int, status string) ([]entity.Volunteer, int, error) {
	var volunteers []entity.Volunteer
	var totalCount int

	// Calculate offset
	offset := (page - 1) * limit

	// Build query with optional status filter
	query := `SELECT * FROM volunteers`
	countQuery := `SELECT COUNT(*) FROM volunteers`
	args := make(map[string]interface{})

	if status != "" {
		query += ` WHERE status = :status`
		countQuery += ` WHERE status = :status`
		args["status"] = status
	}

	query += ` ORDER BY created_at DESC LIMIT :limit OFFSET :offset`
	args["limit"] = limit
	args["offset"] = offset

	// Get total count
	countStmt, err := r.db.PrepareNamed(countQuery)
	if err != nil {
		return nil, 0, err
	}
	defer countStmt.Close()

	if err := countStmt.Get(&totalCount, args); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, 0, err
	}
	defer stmt.Close()

	if err := stmt.Select(&volunteers, args); err != nil {
		return nil, 0, err
	}

	return volunteers, totalCount, nil
}

// GetByID retrieves a volunteer application by its ID
func (r *volunteerRepository) GetByID(id string) (*entity.Volunteer, error) {
	var volunteer entity.Volunteer
	query := `SELECT * FROM volunteers WHERE id = ?`

	err := r.db.Get(&volunteer, query, id)
	if err != nil {
		return nil, err
	}

	return &volunteer, nil
}

// GetByUserID retrieves a volunteer application by user ID
func (r *volunteerRepository) GetByUserID(userID string) (*entity.Volunteer, error) {
	var volunteer entity.Volunteer
	query := `SELECT * FROM volunteers WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`

	err := r.db.Get(&volunteer, query, userID)
	if err != nil {
		return nil, err
	}

	return &volunteer, nil
}

// HasPendingApplication checks if a user has a pending volunteer application
func (r *volunteerRepository) HasPendingApplication(userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM volunteers WHERE user_id = ? AND status = 'pending'`

	err := r.db.Get(&count, query, userID)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateStatus updates the status of a volunteer application
func (r *volunteerRepository) UpdateStatus(id string, status entity.ApplicationStatus) error {
	query := `UPDATE volunteers SET status = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
