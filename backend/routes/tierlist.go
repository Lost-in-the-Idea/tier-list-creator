package routes

import (
	"net/http"

	"tierlist/database"
	"tierlist/middleware"
	"tierlist/models"

	"github.com/gin-gonic/gin"
)
	
func SetupTierlistRoutes (api *gin.RouterGroup, db *database.Database) {
	tierlist := api.Group("/tierlist")
	tierlist.Use(middleware.AuthRequired(db))

	tierlist.GET("/", func(c *gin.Context) { getAllTierlists(c, db) })
	tierlist.GET("/:id", func(c *gin.Context) { getTierlistById(c, db) })
	tierlist.POST("/", func(c *gin.Context) { createTierlist(c, db) })
	tierlist.PUT("/:id", func(c *gin.Context) { updateTierlist(c, db) })
	tierlist.DELETE("/:id", func(c *gin.Context) { deleteTierList(c, db) })

	tierlist.POST("/:id/tier", func(c *gin.Context) { addTier(c, db) })
	tierlist.PUT("/:id/tier/:tierId", func(c *gin.Context) { updateTier(c, db) })
	tierlist.DELETE("/:id/tier/:tierId", func(c *gin.Context) { deleteTier(c, db) })
	tierlist.POST("/:id/item", func(c *gin.Context) { addItem(c, db) })
	tierlist.PUT("/:id/item/:itemId", func(c *gin.Context) { updateItem(c, db) })
	tierlist.DELETE("/:id/item/:itemId", func(c *gin.Context) { deleteItem(c, db) })
}

func getAllTierlists(c *gin.Context, db *database.Database) {
	var tierlists []models.Tierlist
	if err := db.DB.Find(&tierlists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return			
	}
	c.JSON(http.StatusOK, tierlists)
}

func getTierlistById(c *gin.Context, db *database.Database) {
	var tierlist models.Tierlist
	id := c.Param("id")
	if err := db.DB.First(&tierlist, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tierlist not found"})
		return
	}	
	c.JSON(http.StatusOK, tierlist)
}

func createTierlist(c *gin.Context, db *database.Database) {
	var request models.TierlistRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	tierlist := models.Tierlist{
		Name: request.Name,
		Description: request.Description,
		CreatorID: request.Creator,
	}

	if err := db.DB.Create(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	for _, tier := range request.Tiers {
		newTier := models.Tier{
			TierlistID: tierlist.ID,
			Text: tier.Name,
			Colour: tier.Colour,
		}
		if err := db.DB.Create(&newTier).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
			return
		}
	}

	for _, item := range request.Items {
		newItem := models.Item{
			TierlistID: tierlist.ID,
			Text: item.Text,
			Image: item.Image,
			TierText: item.Tier,
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
	if err := db.DB.First(&tierlist, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tierlist not found"})
		return
	}


	tierlist.Name = request.Name
    tierlist.Description = request.Description
    tierlist.CreatorID = request.CreatorID
    tierlist.Version = request.Version


	if err := db.DB.Save(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

c.JSON(http.StatusOK, tierlist)
}

func deleteTierList(c *gin.Context, db *database.Database) {
	id := c.Param("id")
	var tierlist models.Tierlist
	if err := db.DB.First(&tierlist, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tierlist not found"})
		return
	}

	if err := db.DB.Delete(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	if err := db.DB.Where("tierlist_id = ?", id).Delete(&models.Tier{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	if err := db.DB.Where("tierlist_id = ?", id).Delete(&models.Item{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Tierlist deleted"})
}

func addTier(c *gin.Context, db *database.Database) {
	var request models.Tier
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	tier := models.Tier{
		TierlistID: request.TierlistID,
		Text: request.Text,
		Colour: request.Colour,
	}

	if err := db.DB.Create(&tier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusCreated, tier)
}

func updateTier(c *gin.Context, db *database.Database) {
	var request models.Tier
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	tierID := c.Param("tierId")
	var tier models.Tier
	if err := db.DB.First(&tier, tierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tier not found"})
		return
	}

	tier.Text = request.Text
	tier.Colour = request.Colour

	if err := db.DB.Save(&tier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}	

	c.JSON(http.StatusOK, tier)
}

func deleteTier(c *gin.Context, db *database.Database) {
	tierID := c.Param("tierId")
	var tier models.Tier
	if err := db.DB.First(&tier, tierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tier not found"})
		return
	}

	if err := db.DB.Delete(&tier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Tier deleted"})
}

func addItem(c *gin.Context, db *database.Database) {
	var request models.Item
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	item := models.Item{
		TierlistID: request.TierlistID,
		Text: request.Text,
		Image: request.Image,
		TierText: request.TierText,
	}

	if err := db.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func updateItem(c *gin.Context, db *database.Database) {
	var request models.Item
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	itemID := c.Param("itemId")
	var item models.Item
	if err := db.DB.First(&item, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Item not found"})
		return
	}

	item.Text = request.Text
	item.Image = request.Image
	item.TierText = request.TierText

	if err := db.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func deleteItem(c *gin.Context, db *database.Database) {
	itemID := c.Param("itemId")
	var item models.Item
	if err := db.DB.First(&item, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Item not found"})
		return
	}

	if err := db.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Item deleted"})
}