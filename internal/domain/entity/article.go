package entity

import "time"

// Article represents an article entity in the system
type Article struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Slug        string     `db:"slug"`
	Content     string     `db:"content"`
	ImageURL    *string    `db:"image_url"`
	AuthorID    *string    `db:"author_id"`
	PublishedAt *time.Time `db:"published_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
