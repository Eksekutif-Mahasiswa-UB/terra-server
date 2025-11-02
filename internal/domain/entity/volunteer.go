package entity

import "time"

// Gender represents the gender of a volunteer
type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

// VolunteerStatus represents the status/occupation of a volunteer
type VolunteerStatus string

const (
	VolunteerStatusStudent    VolunteerStatus = "Student"
	VolunteerStatusHighSchool VolunteerStatus = "High School"
	VolunteerStatusEmployee   VolunteerStatus = "Employee"
)

// Volunteer represents a volunteer entity in the system
type Volunteer struct {
	ID              string          `db:"id" json:"id"`
	FullName        string          `db:"full_name" json:"full_name"`
	Email           string          `db:"email" json:"email"`
	PhoneNumber     string          `db:"phone_number" json:"phone_number"`
	BirthDate       time.Time       `db:"birth_date" json:"birth_date"`
	Gender          Gender          `db:"gender" json:"gender"`
	Domicile        string          `db:"domicile" json:"domicile"`
	Status          VolunteerStatus `db:"status" json:"status"`
	Interest        string          `db:"interest" json:"interest"`
	CertificateName string          `db:"certificate_name" json:"certificate_name"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}
