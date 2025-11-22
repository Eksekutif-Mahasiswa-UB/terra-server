package handler

import (
	"net/http"
	"strings"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/article/service"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/dto"
	"github.com/gin-gonic/gin"
)

// ArticleHandler handles article-related HTTP requests
type ArticleHandler struct {
	articleService service.ArticleService
}

// NewArticleHandler creates a new instance of ArticleHandler
func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

// CreateArticle handles the create article endpoint
// @Summary Create new article
// @Tags articles
// @Accept json
// @Produce json
// @Param request body dto.CreateArticleRequest true "Create Article Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/articles [post]
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var request dto.CreateArticleRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Additional validation
	if strings.TrimSpace(request.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Title cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Content cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.Category) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Category cannot be empty",
		})
		return
	}

	// Extract userID from context (set by authentication middleware)
	// For now, we'll use a placeholder. In production, this should come from JWT middleware
	authorID, exists := c.Get("userID")
	if !exists {
		// Fallback: if no middleware sets userID, we could reject or use a default
		// For now, rejecting unauthorized requests
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication required",
			"message": "User must be authenticated to create articles",
		})
		return
	}

	authorIDStr, ok := authorID.(string)
	if !ok || strings.TrimSpace(authorIDStr) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication error",
			"message": "Invalid user authentication",
		})
		return
	}

	// Call service to create article
	article, err := h.articleService.CreateArticle(request, authorIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Article creation failed",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	articleResponse := dto.ToArticleResponse(article)

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Article created successfully",
		"data":    articleResponse,
	})
}

// GetAllArticles handles the get all articles endpoint with optional category filter
// @Summary Get all articles
// @Tags articles
// @Produce json
// @Param category query string false "Filter by category"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/articles [get]
func (h *ArticleHandler) GetAllArticles(c *gin.Context) {
	// Get optional category query parameter
	category := c.Query("category")

	// Call service to get articles (with optional filter)
	articles, err := h.articleService.GetAllArticles(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve articles",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response list
	articleResponses := dto.ToArticleResponseList(articles)

	// Return success response
	message := "Articles retrieved successfully"
	if category != "" {
		message = "Articles retrieved successfully for category: " + category
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    articleResponses,
	})
}

// GetArticleBySlug handles the get article by slug endpoint
// @Summary Get article by slug
// @Tags articles
// @Produce json
// @Param slug path string true "Article Slug"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/articles/{slug} [get]
func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	// Get article slug from URL parameter
	slug := c.Param("slug")

	if strings.TrimSpace(slug) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Article slug is required",
		})
		return
	}

	// Call service to get article by slug
	article, err := h.articleService.GetArticleBySlug(slug)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "article not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error":   "Failed to retrieve article",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	articleResponse := dto.ToArticleResponse(article)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Article retrieved successfully",
		"data":    articleResponse,
	})
}

// UpdateArticle handles the update article endpoint
// @Summary Update article
// @Tags articles
// @Accept json
// @Produce json
// @Param slug path string true "Article Slug"
// @Param request body dto.UpdateArticleRequest true "Update Article Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/articles/{slug} [put]
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	// Get article slug from URL parameter
	slug := c.Param("slug")

	if strings.TrimSpace(slug) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Article slug is required",
		})
		return
	}

	var request dto.UpdateArticleRequest

	// Bind JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Additional validation
	if strings.TrimSpace(request.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Title cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Content cannot be empty",
		})
		return
	}

	if strings.TrimSpace(request.Category) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Category cannot be empty",
		})
		return
	}

	// Call service to update article
	article, err := h.articleService.UpdateArticle(slug, request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "article not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error":   "Article update failed",
			"message": err.Error(),
		})
		return
	}

	// Convert to DTO response
	articleResponse := dto.ToArticleResponse(article)

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
		"data":    articleResponse,
	})
}

// DeleteArticle handles the delete article endpoint
// @Summary Delete article
// @Tags articles
// @Produce json
// @Param slug path string true "Article Slug"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/articles/{slug} [delete]
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	// Get article slug from URL parameter
	slug := c.Param("slug")

	if strings.TrimSpace(slug) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation error",
			"message": "Article slug is required",
		})
		return
	}

	// Call service to delete article
	err := h.articleService.DeleteArticle(slug)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "article not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error":   "Article deletion failed",
			"message": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Article deleted successfully",
	})
}
