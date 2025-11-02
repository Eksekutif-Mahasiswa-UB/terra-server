package entity

import "time"

// Article represents an article entity in the system
type Article struct {
	ID          string     `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Slug        string     `db:"slug" json:"slug"`
	Content     string     `db:"content" json:"content"`
	ImageURL    *string    `db:"image_url" json:"image_url,omitempty"`
	AuthorID    *string    `db:"author_id" json:"author_id,omitempty"`
	PublishedAt *time.Time `db:"published_at" json:"published_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}
