package service

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/article/repository"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/google/uuid"
)

// ArticleService defines the interface for article business logic
type ArticleService interface {
	CreateArticle(request dto.CreateArticleRequest, authorID string) (*entity.Article, error)
	GetAllArticles(category string) ([]entity.Article, error)
	GetArticleBySlug(slug string) (*entity.Article, error)
	UpdateArticle(slug string, request dto.UpdateArticleRequest) (*entity.Article, error)
	DeleteArticle(slug string) error
}

// articleService is the concrete implementation of ArticleService
type articleService struct {
	articleRepo repository.ArticleRepository
}

// NewArticleService creates a new instance of ArticleService
func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
	return &articleService{articleRepo: articleRepo}
}

// GenerateSlug converts a title to a URL-friendly slug
func GenerateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove all non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// CreateArticle handles the business logic for creating a new article
func (s *articleService) CreateArticle(request dto.CreateArticleRequest, authorID string) (*entity.Article, error) {
	// Validate required fields
	if strings.TrimSpace(request.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	if strings.TrimSpace(request.Content) == "" {
		return nil, errors.New("content cannot be empty")
	}

	if strings.TrimSpace(request.Category) == "" {
		return nil, errors.New("category cannot be empty")
	}

	if strings.TrimSpace(authorID) == "" {
		return nil, errors.New("author ID is required")
	}

	// Generate slug from title
	slug := GenerateSlug(request.Title)
	if slug == "" {
		return nil, errors.New("unable to generate valid slug from title")
	}

	// Check if slug already exists
	existingArticle, err := s.articleRepo.GetArticleBySlug(slug)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to check slug uniqueness")
	}

	// If slug exists, append UUID suffix
	if existingArticle != nil {
		slug = slug + "-" + uuid.NewString()[:8]
	}

	// Create new article entity
	newArticle := &entity.Article{
		ID:       uuid.NewString(),
		Title:    request.Title,
		Slug:     slug,
		Content:  request.Content,
		ImageURL: request.ImageURL,
		Category: request.Category,
		AuthorID: authorID,
	}

	// Save to database
	if err := s.articleRepo.CreateArticle(newArticle); err != nil {
		return nil, errors.New("failed to create article")
	}

	return newArticle, nil
}

// GetAllArticles retrieves all articles with optional category filter
func (s *articleService) GetAllArticles(category string) ([]entity.Article, error) {
	articles, err := s.articleRepo.GetAllArticles(category)
	if err != nil {
		return nil, errors.New("failed to retrieve articles")
	}

	return articles, nil
}

// GetArticleBySlug retrieves an article by its slug
func (s *articleService) GetArticleBySlug(slug string) (*entity.Article, error) {
	if strings.TrimSpace(slug) == "" {
		return nil, errors.New("slug cannot be empty")
	}

	article, err := s.articleRepo.GetArticleBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("article not found")
		}
		return nil, errors.New("failed to retrieve article")
	}

	return article, nil
}

// UpdateArticle handles the business logic for updating an article
func (s *articleService) UpdateArticle(slug string, request dto.UpdateArticleRequest) (*entity.Article, error) {
	if strings.TrimSpace(slug) == "" {
		return nil, errors.New("slug cannot be empty")
	}

	// Validate required fields
	if strings.TrimSpace(request.Title) == "" {
		return nil, errors.New("title cannot be empty")
	}

	if strings.TrimSpace(request.Content) == "" {
		return nil, errors.New("content cannot be empty")
	}

	if strings.TrimSpace(request.Category) == "" {
		return nil, errors.New("category cannot be empty")
	}

	// Check if article exists
	existingArticle, err := s.articleRepo.GetArticleBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("article not found")
		}
		return nil, errors.New("failed to retrieve article")
	}

	// Generate new slug if title changed
	newSlug := slug
	if request.Title != existingArticle.Title {
		newSlug = GenerateSlug(request.Title)
		if newSlug == "" {
			return nil, errors.New("unable to generate valid slug from title")
		}

		// Check if new slug already exists (and it's not the current article)
		if newSlug != slug {
			checkArticle, err := s.articleRepo.GetArticleBySlug(newSlug)
			if err != nil && err != sql.ErrNoRows {
				return nil, errors.New("failed to check slug uniqueness")
			}
			if checkArticle != nil {
				newSlug = newSlug + "-" + uuid.NewString()[:8]
			}
		}
	}

	// Update article fields
	existingArticle.Title = request.Title
	existingArticle.Slug = newSlug
	existingArticle.Content = request.Content
	existingArticle.ImageURL = request.ImageURL
	existingArticle.Category = request.Category

	// Save updates to database
	if err := s.articleRepo.UpdateArticle(existingArticle); err != nil {
		return nil, errors.New("failed to update article")
	}

	return existingArticle, nil
}

// DeleteArticle handles the business logic for deleting an article
func (s *articleService) DeleteArticle(slug string) error {
	if strings.TrimSpace(slug) == "" {
		return errors.New("slug cannot be empty")
	}

	// Check if article exists
	_, err := s.articleRepo.GetArticleBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("article not found")
		}
		return errors.New("failed to retrieve article")
	}

	// Delete article
	if err := s.articleRepo.DeleteArticle(slug); err != nil {
		return errors.New("failed to delete article")
	}

	return nil
}
