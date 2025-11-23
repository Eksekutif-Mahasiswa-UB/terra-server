package repository

import (
	"database/sql"
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// DonationRepository defines the interface for donation data operations
type DonationRepository interface {
	Create(donation *entity.Donation) error
	GetByUserID(userID string) ([]entity.Donation, error)
	GetAll() ([]entity.Donation, error)
	GetByID(id string) (*entity.Donation, error)
	UpdateStatus(id string, status entity.DonationStatus) error
}

// donationRepository is the concrete implementation of DonationRepository
type donationRepository struct {
	db *sqlx.DB
}

// NewDonationRepository creates a new instance of DonationRepository
func NewDonationRepository(db *sqlx.DB) DonationRepository {
	return &donationRepository{db: db}
}

// Create inserts a new donation into the database
func (r *donationRepository) Create(donation *entity.Donation) error {
	donation.CreatedAt = time.Now()
	donation.UpdatedAt = time.Now()

	query := `INSERT INTO donations (id, user_id, program_id, amount, payment_method, status, proof_image, created_at, updated_at) 
			  VALUES (:id, :user_id, :program_id, :amount, :payment_method, :status, :proof_image, :created_at, :updated_at)`

	_, err := r.db.NamedExec(query, donation)
	return err
}

// GetByUserID retrieves all donations made by a specific user
func (r *donationRepository) GetByUserID(userID string) ([]entity.Donation, error) {
	var donations []entity.Donation
	query := `SELECT * FROM donations WHERE user_id = ? ORDER BY created_at DESC`

	err := r.db.Select(&donations, query, userID)
	if err != nil {
		return nil, err
	}

	return donations, nil
}

// GetAll retrieves all donations from the database
func (r *donationRepository) GetAll() ([]entity.Donation, error) {
	var donations []entity.Donation
	query := `SELECT * FROM donations ORDER BY created_at DESC`

	err := r.db.Select(&donations, query)
	if err != nil {
		return nil, err
	}

	return donations, nil
}

// GetByID retrieves a donation by its ID
func (r *donationRepository) GetByID(id string) (*entity.Donation, error) {
	var donation entity.Donation
	query := `SELECT * FROM donations WHERE id = ?`

	err := r.db.Get(&donation, query, id)
	if err != nil {
		return nil, err
	}

	return &donation, nil
}

// UpdateStatus updates the status of a donation
func (r *donationRepository) UpdateStatus(id string, status entity.DonationStatus) error {
	query := `UPDATE donations SET status = ?, updated_at = ? WHERE id = ?`

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
