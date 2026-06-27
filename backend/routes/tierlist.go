package routes

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tierlist/database/models"
	"tierlist/dto"
	"tierlist/services"
)

func SetupTierlistRoutes(api *gin.RouterGroup, svc *services.TierlistService, authRequired, optionalAuth gin.HandlerFunc) {
	tierlists := api.Group("/tierlists")
	tierlists.GET("/:id/results", func(c *gin.Context) { getTierlistResults(c, svc) })
	tierlists.GET("/:id", optionalAuth, func(c *gin.Context) { getTierlistById(c, svc) })
	tierlists.POST("/", authRequired, func(c *gin.Context) { createNewTierlist(c, svc) })
	tierlists.POST("/:id/submit", authRequired, func(c *gin.Context) { submitTierlist(c, svc) })
	tierlists.DELETE("/:id", authRequired, func(c *gin.Context) { deleteById(c, svc) })
}

func getTierlistById(c *gin.Context, svc *services.TierlistService) {
	id := c.Param("id")
	var userID *uuid.UUID
	if u, exists := c.Get("user"); exists {
		uid := u.(models.User).ID
		userID = &uid
	}
	result, err := svc.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tierlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func createNewTierlist(c *gin.Context, svc *services.TierlistService) {
	var req dto.CreateTierlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creatorID := c.MustGet("user").(models.User).ID
	result, err := svc.Create(req, creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func submitTierlist(c *gin.Context, svc *services.TierlistService) {
	id := c.Param("id")
	var req dto.SubmitRankingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.MustGet("user").(models.User).ID
	err := svc.Submit(id, userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Tierlist not found"})
		case errors.Is(err, services.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": "Already submitted"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		}
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Submission successful"})
}

func getTierlistResults(c *gin.Context, svc *services.TierlistService) {
	id := c.Param("id")
	result, err := svc.GetResults(id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tierlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func deleteById(c *gin.Context, svc *services.TierlistService) {
	id := c.Param("id")
	userID := c.MustGet("user").(models.User).ID
	err := svc.Delete(id, userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Tierlist not found"})
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		}
		return
	}
	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Tierlist deleted successfully"})
}
