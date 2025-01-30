package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tierlist/database"
	"tierlist/models"
)

func SetupTierlistRoutes (r *gin.Engine) {
	tierlist := r.Group("/tierlist")

	tierlist.GET("/", getAllTierlists)
	tierlist.GET("/:id", getTierlistById)
	tierlist.POST("/", createTierlist)
	tierlist.PUT("/:id", updateTierlist)
	tierlist.DELETE("/:id", deleteTierList)

	tierlist.POST("/:id/tier", addTier)
	tierlist.PUT("/:id/tier/:tierId", updateTier)
	tierlist.DELETE("/:id/tier/:tierId", deleteTier)
	tierlist.POST("/:id/item", addItem)
	tierlist.PUT("/:id/item/:itemId", updateItem)
	tierlist.DELETE("/:id/item/:itemId", deleteItem)
}

func getAllTierlists(c *gin.Context) {
	var tierlists []models.Tierlist
	if err := database.DB.Find(&tierlists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}
	c.JSON(http.StatusOK, tierlists)
}

func getTierlistById(c *gin.Context) {
	var tierlists []models.Tierlist
	id := c.Param("id")
	if err := database.DB.First(&tierlists, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tierlist not found"})
		return
	}
	c.JSON(http.StatusOK, tierlists)
}

func createTierlist(c *gin.Context) {
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

	if err := database.DB.Create(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	for _, tier := range request.Tiers {
		newTier := models.Tier{
			TierlistID: tierlist.ID,
			Text: tier.Name,
			Colour: tier.Colour,
		}
		if err := database.DB.Create(&newTier).Error; err != nil {
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
		if err := database.DB.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
			return
		}
	}

	c.JSON(http.StatusCreated, tierlist)
}

func updateTierlist(c *gin.Context) {
	var request models.Tierlist
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	id := c.Param("id")
	var tierlist models.Tierlist
	if err := database.DB.First(&tierlist, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tierlist not found"})
		return
	}


	tierlist.Name = request.Name
    tierlist.Description = request.Description
    tierlist.CreatorID = request.CreatorID
    tierlist.Version = request.Version


	if err := database.DB.Save(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

c.JSON(http.StatusOK, tierlist)
}

func deleteTierList(c *gin.Context) {
	id := c.Param("id")
	var tierlist models.Tierlist
	if err := database.DB.First(&tierlist, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tierlist not found"})
		return
	}

	if err := database.DB.Delete(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	if err := database.DB.Where("tierlist_id = ?", id).Delete(&models.Tier{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	if err := database.DB.Where("tierlist_id = ?", id).Delete(&models.Item{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Tierlist deleted"})
}

func addTier(c *gin.Context) {
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

	if err := database.DB.Create(&tier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusCreated, tier)
}

func updateTier(c *gin.Context) {
	var request models.Tier
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	tierID := c.Param("tierId")
	var tier models.Tier
	if err := database.DB.First(&tier, tierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tier not found"})
		return
	}

	tier.Text = request.Text
	tier.Colour = request.Colour

	if err := database.DB.Save(&tier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}	

	c.JSON(http.StatusOK, tier)
}

func deleteTier(c *gin.Context) {
	tierID := c.Param("tierId")
	var tier models.Tier
	if err := database.DB.First(&tier, tierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Tier not found"})
		return
	}

	if err := database.DB.Delete(&tier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Tier deleted"})
}

func addItem(c *gin.Context) {
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

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func updateItem(c *gin.Context) {
	var request models.Item
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	itemID := c.Param("itemId")
	var item models.Item
	if err := database.DB.First(&item, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Item not found"})
		return
	}

	item.Text = request.Text
	item.Image = request.Image
	item.TierText = request.TierText

	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func deleteItem(c *gin.Context) {
	itemID := c.Param("itemId")
	var item models.Item
	if err := database.DB.First(&item, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Item not found"})
		return
	}

	if err := database.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Item deleted"})
}