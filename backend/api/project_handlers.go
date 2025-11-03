package api

import (
	"ciphersafe/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	DB *gorm.DB
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{DB: db}
}

type projectInput struct {
	Name string `json:"name" binding:"required"`
}

// CreateProject handles creation of a new project
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var input projectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := getUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	project := models.Project{
		Name:    input.Name,
		OwnerID: userID,
	}

	if err := h.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetProjects lists all projects for the authenticated user
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID, exists := getUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var projects []models.Project
	// Eager load secrets to show them
	if err := h.DB.Preload("Secrets").Where("owner_id = ?", userID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}
