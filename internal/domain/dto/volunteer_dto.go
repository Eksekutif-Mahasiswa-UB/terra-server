package dto

import (
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
)

// VolunteerResponse represents the volunteer data transfer object for API responses
type VolunteerResponse struct {
	ID              string                 `json:"id"`
	FullName        string                 `json:"full_name"`
	Email           string                 `json:"email"`
	PhoneNumber     string                 `json:"phone_number"`
	BirthDate       time.Time              `json:"birth_date"`
	Gender          entity.Gender          `json:"gender"`
	Domicile        string                 `json:"domicile"`
	Status          entity.VolunteerStatus `json:"status"`
	Interest        string                 `json:"interest"`
	CertificateName string                 `json:"certificate_name"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
