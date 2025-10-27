package model

import "time"


type Article struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Slug        string    `json:"slug" db:"slug"`
	Content     string    `json:"content" db:"content"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	AuthorID    *string   `json:"author_id" db:"author_id"` 
	PublishedAt *time.Time `json:"published_at" db:"published_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}