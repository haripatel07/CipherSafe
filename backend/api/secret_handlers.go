package api

import (
	"ciphersafe/models"
	"ciphersafe/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SecretHandler struct {
	DB *gorm.DB
}

func NewSecretHandler(db *gorm.DB) *SecretHandler {
	return &SecretHandler{DB: db}
}

type secretInput struct {
	ProjectID uint   `json:"project_id" binding:"required"`
	Key       string `json:"key" binding:"required"`
	Value     string `json:"value" binding:"required"` // This is the PLAINTEXT value
}

// DecryptedSecret is a struct for sending secrets to the user
type DecryptedSecret struct {
	ID        uint   `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"` // This will hold the DECRYPTED value
	ProjectID uint   `json:"project_id"`
}

// verifyProjectOwnership is a crucial helper function
func verifyProjectOwnership(db *gorm.DB, userID, projectID uint) bool {
	var project models.Project
	if err := db.First(&project, projectID).Error; err != nil {
		return false // Project not found
	}
	return project.OwnerID == userID
}

// CreateSecret encrypts and saves a new secret
func (h *SecretHandler) CreateSecret(c *gin.Context) {
	var input secretInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := getUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify the authenticated user actually owns the project they're adding a secret to
	if !verifyProjectOwnership(h.DB, userID, input.ProjectID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission for this project"})
		return
	}

	// Encrypt the secret value before saving
	encryptedValue, err := services.Encrypt(input.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt secret"})
		return
	}

	secret := models.Secret{
		ProjectID: input.ProjectID,
		Key:       input.Key,
		Value:     encryptedValue, // Save the ENCRYPTED value
	}

	if err := h.DB.Create(&secret).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save secret"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Secret created successfully"})
}

// GetSecretsForProject decrypts and returns all secrets for a project
func (h *SecretHandler) GetSecretsForProject(c *gin.Context) {
	projectIDStr := c.Param("projectID")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	userID, exists := getUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// *** CRITICAL SECURITY CHECK ***
	if !verifyProjectOwnership(h.DB, userID, uint(projectID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission for this project"})
		return
	}

	var secrets []models.Secret
	if err := h.DB.Where("project_id = ?", projectID).Find(&secrets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve secrets"})
		return
	}

	// Decrypt secrets before sending them
	var decryptedSecrets []DecryptedSecret
	for _, secret := range secrets {
		decryptedValue, err := services.Decrypt(secret.Value)
		if err != nil {
			continue
		}

		decryptedSecrets = append(decryptedSecrets, DecryptedSecret{
			ID:        secret.ID,
			Key:       secret.Key,
			Value:     decryptedValue, // Send the DECRYPTED value
			ProjectID: secret.ProjectID,
		})
	}

	c.JSON(http.StatusOK, decryptedSecrets)
}

// DeleteSecret deletes a specific secret
func (h *SecretHandler) DeleteSecret(c *gin.Context) {
	secretIDStr := c.Param("secretID")
	secretID, err := strconv.ParseUint(secretIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid secret ID"})
		return
	}

	userID, exists := getUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Verify the user owns the project that this secret belongs to.
	var secret models.Secret
	if err := h.DB.First(&secret, uint(secretID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !verifyProjectOwnership(h.DB, userID, secret.ProjectID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission for this secret"})
		return
	}

	// All checks passed, delete the secret
	if err := h.DB.Delete(&secret).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete secret"})
		return
	}

	c.JSON(http.StatusNoContent, nil) // 204 No Content is standard for successful delete
}
