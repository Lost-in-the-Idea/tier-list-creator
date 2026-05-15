package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tierlist/database"
	"tierlist/database/models"
	"tierlist/dto"
	"tierlist/middleware"
)

func SetupTierlistRoutes(api *gin.RouterGroup, db *database.Database) {
	tierlist := api.Group("/tierlist")
	tierlist.Use(middleware.AuthRequired(db))

	tierlist.GET("/", func(c *gin.Context) { getAllTierlists(c, db) })
	tierlist.GET("/:id", func(c *gin.Context) { getTierlistById(c, db) })
	tierlist.POST("/", func(c *gin.Context) { createTierlist(c, db) })
	tierlist.PUT("/:id", func(c *gin.Context) { updateTierlist(c, db) })
	tierlist.DELETE("/:id", func(c *gin.Context) { deleteTierList(c, db) })

	tierlist.POST("/:id/item", func(c *gin.Context) { addItem(c, db) })
	tierlist.PUT("/:id/item/:itemId", func(c *gin.Context) { updateItem(c, db) })
	tierlist.DELETE("/:id/item/:itemId", func(c *gin.Context) { deleteItem(c, db) })
}

func getAllTierlists(c *gin.Context, db *database.Database) {
	var tierlists []models.Tierlist
	if err := db.DB.Preload("TierlistItems").Find(&tierlists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, tierlists)
}

func getTierlistById(c *gin.Context, db *database.Database) {
	id := c.Param("id")
	var tierlist models.Tierlist
	if err := db.DB.Where("id = ?", id).Preload("TierlistItems").First(&tierlist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tierlist not found"})
		return
	}
	c.JSON(http.StatusOK, tierlist)
}

func createTierlist(c *gin.Context, db *database.Database) {
	var request dto.CreateTierlistRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	user, _ := c.Get("user")
	creator := user.(models.User)

	tierlist := models.Tierlist{
		Title:       request.Title,
		Description: request.Description,
		CreatorID:   creator.ID,
		ShareCode:   uuid.New().String(),
		ExpiryTime:  request.ExpiryTime,
	}

	if err := db.DB.Create(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	for _, item := range request.Items {
		newItem := models.TierlistItem{
			TierlistID: tierlist.ID,
			Name:       item.Name,
			ImageURL:   item.ImageURL,
			SortOrder:  item.SortOrder,
		}
		if err := db.DB.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
			return
		}
	}

	c.JSON(http.StatusCreated, tierlist)
}

func updateTierlist(c *gin.Context, db *database.Database) {
	var request models.Tierlist
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	id := c.Param("id")
	var tierlist models.Tierlist
	if err := db.DB.Where("id = ?", id).First(&tierlist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tierlist not found"})
		return
	}

	tierlist.Title = request.Title
	tierlist.Description = request.Description

	if err := db.DB.Save(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, tierlist)
}

func deleteTierList(c *gin.Context, db *database.Database) {
	id := c.Param("id")
	var tierlist models.Tierlist
	if err := db.DB.Where("id = ?", id).First(&tierlist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tierlist not found"})
		return
	}

	if err := db.DB.Delete(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Tierlist deleted"})
}

func addItem(c *gin.Context, db *database.Database) {
	var request dto.CreateItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	tierlistID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid tierlist ID"})
		return
	}

	item := models.TierlistItem{
		TierlistID: tierlistID,
		Name:       request.Name,
		ImageURL:   request.ImageURL,
		SortOrder:  request.SortOrder,
	}

	if err := db.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func updateItem(c *gin.Context, db *database.Database) {
	var request dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	itemID := c.Param("itemId")
	var item models.TierlistItem
	if err := db.DB.Where("id = ?", itemID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Item not found"})
		return
	}

	item.Name = request.Name
	item.ImageURL = request.ImageURL
	item.SortOrder = request.SortOrder

	if err := db.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func deleteItem(c *gin.Context, db *database.Database) {
	itemID := c.Param("itemId")
	var item models.TierlistItem
	if err := db.DB.Where("id = ?", itemID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Item not found"})
		return
	}

	if err := db.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Item deleted"})
}
