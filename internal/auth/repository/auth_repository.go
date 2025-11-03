package repository

import (
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// AuthRepository defines the interface for auth data operations
type AuthRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	CreateRefreshToken(token *entity.RefreshToken) error
}

// authRepository is the concrete implementation of AuthRepository
type authRepository struct {
	db *sqlx.DB
}

// NewAuthRepository creates a new instance of AuthRepository
func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

// CreateUser inserts a new user into the database
func (r *authRepository) CreateUser(user *entity.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `INSERT INTO users (id, full_name, email, password, role, auth_method, created_at, updated_at)
          VALUES (:id, :full_name, :email, :password, :role, :auth_method, :created_at, :updated_at)`

	_, err := r.db.NamedExec(query, user)
	return err
}

// GetUserByEmail retrieves a user by their email address
func (r *authRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	query := `SELECT * FROM users WHERE email = ?`

	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID
func (r *authRepository) GetUserByID(id string) (*entity.User, error) {
	var user entity.User
	query := `SELECT * FROM users WHERE id = ?`

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateRefreshToken inserts a new refresh token into the database
func (r *authRepository) CreateRefreshToken(token *entity.RefreshToken) error {
	token.CreatedAt = time.Now()

	query := `INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at) 
			  VALUES (:id, :user_id, :token, :expires_at, :created_at)`

	_, err := r.db.NamedExec(query, token)
	return err
}
