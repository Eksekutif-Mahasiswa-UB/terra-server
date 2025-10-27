package model

import "time"


type User struct {
    ID          string    `json:"id" db:"id"`
    NamaLengkap string    `json:"nama_lengkap" db:"nama_lengkap"`
    Email       string    `json:"email" db:"email"`
    Password    string    `json:"-" db:"password"` 
    Role        string    `json:"role" db:"role"`
    GoogleID    *string   `json:"-" db:"google_id"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}