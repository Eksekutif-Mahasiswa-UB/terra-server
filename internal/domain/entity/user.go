package entity

import "time"

// User represents a user entity in the system
type User struct {
	ID         string    `db:"id" json:"id"`
	FullName   string    `db:"nama_lengkap" json:"full_name"`
	Email      string    `db:"email" json:"email"`
	Password   string    `db:"password" json:"-"`
	Role       string    `db:"role" json:"role"`
	AuthMethod string    `db:"auth_method" json:"auth_method"`
	GoogleID   *string   `db:"google_id" json:"google_id,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
