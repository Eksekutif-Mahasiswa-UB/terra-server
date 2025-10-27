package model

import "time"


type Volunteer struct {
	ID               string    `json:"id" db:"id"`
	NamaLengkap      string    `json:"nama_lengkap" db:"nama_lengkap"`
	Email            string    `json:"email" db:"email"`
	NoWhatsapp       string    `json:"no_whatsapp" db:"no_whatsapp"`
	TanggalLahir     string    `json:"tanggal_lahir" db:"tanggal_lahir"` 
	JenisKelamin     string    `json:"jenis_kelamin" db:"jenis_kelamin"`
	Domisili         string    `json:"domisili" db:"domisili"`
	Status           string    `json:"status" db:"status"`
	MinatLingkungan  string    `json:"minat_lingkungan" db:"minat_lingkungan"`
	NamaSertifikat   string    `json:"nama_sertifikat" db:"nama_sertifikat"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}