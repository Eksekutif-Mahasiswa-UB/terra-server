package entity

import "time"

// Gender represents the gender of a volunteer
type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

// ApplicationStatus represents the application status of a volunteer
type ApplicationStatus string

const (
	ApplicationStatusPending  ApplicationStatus = "pending"
	ApplicationStatusApproved ApplicationStatus = "approved"
	ApplicationStatusRejected ApplicationStatus = "rejected"
)

// Volunteer represents a volunteer application entity in the system
type Volunteer struct {
	ID          string            `db:"id" json:"id"`
	UserID      string            `db:"user_id" json:"user_id"`
	FullName    string            `db:"full_name" json:"full_name"`
	Email       string            `db:"email" json:"email"`
	Phone       string            `db:"phone" json:"phone"`
	DateOfBirth time.Time         `db:"date_of_birth" json:"date_of_birth"`
	Gender      Gender            `db:"gender" json:"gender"`
	City        string            `db:"city" json:"city"`
	Occupation  string            `db:"occupation" json:"occupation"`
	Interests   string            `db:"interests" json:"interests"`
	Experience  string            `db:"experience" json:"experience"`
	Status      ApplicationStatus `db:"status" json:"status"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `db:"updated_at" json:"updated_at"`
}
