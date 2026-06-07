package routes

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tierlist/database"
	"tierlist/database/models"
	"tierlist/dto"
	"tierlist/middleware"
)

func SetupTierlistRoutes(api *gin.RouterGroup, db *database.Database) {
	tierlists := api.Group("/tierlists")
	tierlists.GET("/:id/results", func(c *gin.Context) { getTierlistResults(c, db) })
	tierlists.GET("/:id", middleware.OptionalAuth(db), func(c *gin.Context) { getTierlistById(c, db) })
	tierlists.POST("/", middleware.AuthRequired(db), func(c *gin.Context) { createNewTierlist(c, db) })
	tierlists.POST("/:id/submit", middleware.AuthRequired(db), func(c *gin.Context) { submitTierlist(c, db) })
	tierlists.DELETE("/:id", middleware.AuthRequired(db), func(c *gin.Context) { deleteById(c, db) })
}

func getTierlistById(c *gin.Context, db *database.Database) {
	id := c.Param("id")
	var tierlist models.Tierlist
	if err := db.DB.Where("id = ?", id).Preload("Creator").Preload("TierlistItems").First(&tierlist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tierlist not found"})
		return
	}

	hasSubmitted := false
	if u, exists := c.Get("user"); exists {
		currentUser := u.(models.User)
		var submission models.Submissions
		if err := db.DB.Where("tierlist_id = ? AND user_id = ?", tierlist.ID, currentUser.ID).First(&submission).Error; err == nil {
			hasSubmitted = true
		}
	}

	items := make([]dto.TierlistItemResponse, len(tierlist.TierlistItems))
	for i, item := range tierlist.TierlistItems {
		items[i] = dto.TierlistItemResponse{
			ID:        item.ID.String(),
			Name:      item.Name,
			ImageURL:  item.ImageURL,
			SortOrder: item.SortOrder,
		}
	}

	c.JSON(http.StatusOK, dto.TierlistResponse{
		ID:          tierlist.ID.String(),
		ShareCode:   tierlist.ShareCode,
		Title:       tierlist.Title,
		Description: tierlist.Description,
		ExpiresAt:   tierlist.ExpiryTime,
		Creator: dto.UserResponse{
			ID:        tierlist.Creator.ID.String(),
			DiscordID: tierlist.Creator.DiscordID,
			Username:  tierlist.Creator.Username,
			Avatar:    tierlist.Creator.Avatar,
		},
		Items:        items,
		HasSubmitted: hasSubmitted,
	})
}

func createNewTierlist(c *gin.Context, db *database.Database) {
	var request dto.CreateTierlistRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, _ := c.Get("user")
	creator := u.(models.User)

	tierlist := models.Tierlist{
		Title:       request.Title,
		Description: request.Description,
		CreatorID:   creator.ID,
		ShareCode:   uuid.New().String(),
		ExpiryTime:  request.ExpiryTime,
	}

	if err := db.DB.Create(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			return
		}
	}

	c.JSON(http.StatusCreated, dto.CreateTierlistResponse{
		ShareCode: tierlist.ShareCode,
		ExpiresAt: tierlist.ExpiryTime,
	})
}

func submitTierlist(c *gin.Context, db *database.Database) {
	id := c.Param("id")
	var tierlist models.Tierlist
	if err := db.DB.Where("id = ?", id).Preload("TierlistItems").First(&tierlist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tierlist not found"})
		return
	}

	u, _ := c.Get("user")
	currentUser := u.(models.User)

	var existing models.Submissions
	if err := db.DB.Where("tierlist_id = ? AND user_id = ?", tierlist.ID, currentUser.ID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already submitted"})
		return
	}

	var request dto.SubmitRankingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validItemIDs := make(map[uuid.UUID]bool, len(tierlist.TierlistItems))
	for _, item := range tierlist.TierlistItems {
		validItemIDs[item.ID] = true
	}

	submission := models.Submissions{
		TierlistID: tierlist.ID,
		UserID:     currentUser.ID,
	}
	if err := db.DB.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}

	for _, ranking := range request.Rankings {
		itemID, err := uuid.Parse(ranking.ItemID)
		if err != nil || !validItemIDs[itemID] {
			db.DB.Delete(&submission)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID: " + ranking.ItemID})
			return
		}
		rankingRecord := models.SubmissionRankings{
			SubmissionID: submission.ID,
			ItemID:       itemID,
			Tier:         ranking.Tier,
		}
		if err := db.DB.Create(&rankingRecord).Error; err != nil {
			db.DB.Delete(&submission)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
			return
		}
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Submission successful"})
}

func getTierlistResults(c *gin.Context, db *database.Database) {
	id := c.Param("id")
	var tierlist models.Tierlist
	if err := db.DB.Where("id = ?", id).Preload("Creator").Preload("TierlistItems").Preload("Submissions.Rankings").First(&tierlist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tierlist not found"})
		return
	}

	tierScores := map[string]int{"S": 6, "A": 5, "B": 4, "C": 3, "D": 2, "F": 1}
	tierPrecedence := []string{"S", "A", "B", "C", "D", "F"}

	type itemStats struct {
		item   models.TierlistItem
		counts map[string]int
		total  int
	}

	statsMap := make(map[uuid.UUID]*itemStats, len(tierlist.TierlistItems))
	for _, item := range tierlist.TierlistItems {
		statsMap[item.ID] = &itemStats{item: item, counts: make(map[string]int)}
	}

	for _, submission := range tierlist.Submissions {
		for _, ranking := range submission.Rankings {
			if s, ok := statsMap[ranking.ItemID]; ok {
				s.counts[ranking.Tier]++
				s.total++
			}
		}
	}

	results := make([]dto.TierResult, 0, len(tierlist.TierlistItems))
	for _, item := range tierlist.TierlistItems {
		s := statsMap[item.ID]

		totalScore := 0
		for tier, count := range s.counts {
			totalScore += tierScores[tier] * count
		}
		var avgScore float64
		if s.total > 0 {
			avgScore = float64(totalScore) / float64(s.total)
		}

		topTier := ""
		topCount := 0
		for _, tier := range tierPrecedence {
			if s.counts[tier] > topCount {
				topCount = s.counts[tier]
				topTier = tier
			}
		}

		results = append(results, dto.TierResult{
			ItemID:       item.ID.String(),
			ItemName:     item.Name,
			ImageURL:     item.ImageURL,
			Counts:       s.counts,
			Total:        s.total,
			TopTier:      topTier,
			AverageScore: avgScore,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].AverageScore > results[j].AverageScore
	})
	for i := range results {
		results[i].Rank = i + 1
	}

	items := make([]dto.TierlistItemResponse, len(tierlist.TierlistItems))
	for i, item := range tierlist.TierlistItems {
		items[i] = dto.TierlistItemResponse{
			ID:        item.ID.String(),
			Name:      item.Name,
			ImageURL:  item.ImageURL,
			SortOrder: item.SortOrder,
		}
	}

	c.JSON(http.StatusOK, dto.TierlistResultResponse{
		Tierlist: dto.TierlistResponse{
			ID:          tierlist.ID.String(),
			ShareCode:   tierlist.ShareCode,
			Title:       tierlist.Title,
			Description: tierlist.Description,
			ExpiresAt:   tierlist.ExpiryTime,
			Creator: dto.UserResponse{
				ID:        tierlist.Creator.ID.String(),
				DiscordID: tierlist.Creator.DiscordID,
				Username:  tierlist.Creator.Username,
				Avatar:    tierlist.Creator.Avatar,
			},
			Items:        items,
			HasSubmitted: false, // results endpoint doesn't need this field, but required by TierlistResponse struct, so just set to false for now
		},
		TotalSubmissions: len(tierlist.Submissions),
		Results:          results,
	})
}

func deleteById(c *gin.Context, db *database.Database) {
	id := c.Param("id")
	var tierlist models.Tierlist
	if err := db.DB.Where("id = ?", id).First(&tierlist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tierlist not found"})
		return
	}

	u, _ := c.Get("user")
	currentUser := u.(models.User)

	if tierlist.CreatorID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := db.DB.Delete(&tierlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Tierlist deleted successfully"})
}
