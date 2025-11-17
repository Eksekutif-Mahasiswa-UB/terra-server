package dto

import "time"

// ArticleResponse represents the article data transfer object for API responses
type ArticleResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	ImageURL    *string    `json:"image_url,omitempty"`
	AuthorID    *string    `json:"author_id,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
