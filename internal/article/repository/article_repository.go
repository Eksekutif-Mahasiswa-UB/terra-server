package repository

import (
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// ArticleRepository defines the interface for article data operations
type ArticleRepository interface {
	CreateArticle(article *entity.Article) error
	GetArticleBySlug(slug string) (*entity.Article, error)
	GetAllArticles(category string) ([]entity.Article, error)
	UpdateArticle(article *entity.Article) error
	DeleteArticle(slug string) error
}

// articleRepository is the concrete implementation of ArticleRepository
type articleRepository struct {
	db *sqlx.DB
}

// NewArticleRepository creates a new instance of ArticleRepository
func NewArticleRepository(db *sqlx.DB) ArticleRepository {
	return &articleRepository{db: db}
}

// CreateArticle inserts a new article into the database
func (r *articleRepository) CreateArticle(article *entity.Article) error {
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	query := `INSERT INTO articles (id, title, slug, content, image_url, category, author_id, published_at, created_at, updated_at) 
			  VALUES (:id, :title, :slug, :content, :image_url, :category, :author_id, :published_at, :created_at, :updated_at)`

	_, err := r.db.NamedExec(query, article)
	return err
}

// GetArticleBySlug retrieves an article by its slug
func (r *articleRepository) GetArticleBySlug(slug string) (*entity.Article, error) {
	var article entity.Article
	query := `SELECT * FROM articles WHERE slug = ?`

	err := r.db.Get(&article, query, slug)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

// GetAllArticles retrieves all articles with optional category filter
func (r *articleRepository) GetAllArticles(category string) ([]entity.Article, error) {
	var articles []entity.Article
	var query string
	var err error

	if category != "" {
		// Filter by category
		query = `SELECT * FROM articles WHERE category = ? ORDER BY created_at DESC`
		err = r.db.Select(&articles, query, category)
	} else {
		// Get all articles
		query = `SELECT * FROM articles ORDER BY created_at DESC`
		err = r.db.Select(&articles, query)
	}

	if err != nil {
		return nil, err
	}

	return articles, nil
}

// UpdateArticle updates an existing article in the database
func (r *articleRepository) UpdateArticle(article *entity.Article) error {
	article.UpdatedAt = time.Now()

	query := `UPDATE articles 
			  SET title = :title, slug = :slug, content = :content, image_url = :image_url, category = :category, updated_at = :updated_at 
			  WHERE slug = :slug`

	_, err := r.db.NamedExec(query, article)
	return err
}

// DeleteArticle deletes an article from the database by slug
func (r *articleRepository) DeleteArticle(slug string) error {
	query := `DELETE FROM articles WHERE slug = ?`
	_, err := r.db.Exec(query, slug)
	return err
}
