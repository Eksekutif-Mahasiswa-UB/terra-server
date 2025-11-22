package dto

import "time"

// CreateArticleRequest represents the request to create a new article
type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ImageURL string `json:"image_url"`
	Category string `json:"category" binding:"required"`
}

// UpdateArticleRequest represents the request to update an article
type UpdateArticleRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ImageURL string `json:"image_url"`
	Category string `json:"category" binding:"required"`
}

// ArticleResponse represents the article data transfer object for API responses
type ArticleResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	ImageURL    string     `json:"image_url"`
	Category    string     `json:"category"`
	AuthorID    string     `json:"author_id"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
