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
	ID              string          `db:"id"`
	FullName        string          `db:"full_name"`
	Email           string          `db:"email"`
	PhoneNumber     string          `db:"phone_number"`
	BirthDate       time.Time       `db:"birth_date"`
	Gender          Gender          `db:"gender"`
	Domicile        string          `db:"domicile"`
	Status          VolunteerStatus `db:"status"`
	Interest        string          `db:"interest"`
	CertificateName string          `db:"certificate_name"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}
