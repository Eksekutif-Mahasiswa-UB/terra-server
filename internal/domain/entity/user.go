package entity

import "time"

// User represents a user entity in the system
type User struct {
	ID         string    `db:"id"`
	FullName   string    `db:"nama_lengkap"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Role       string    `db:"role"`
	AuthMethod string    `db:"auth_method"`
	GoogleID   *string   `db:"google_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
