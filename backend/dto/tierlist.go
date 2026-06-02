package dto

import "time"

type CreateTierlistRequest struct {
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description"`
	ExpiryTime  time.Time           `json:"expiry_time" binding:"required"`
	Items []CreateItemRequest 		`json:"tierlist_items" binding:"required,min=2,max=15,dive,required"` // number of items should be between 2 and 15, and each item is required
}

type CreateItemRequest struct {
	Name      string `json:"name" binding:"required"`
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`
}

type SubmitRankingRequest struct {
	Rankings []RankingRequest `json:"rankings" binding:"required,min=1,dive,required"` // number of rankings should be at least 1 and match the number of items in the tierlist
}

type RankingRequest struct {
	ItemID string `json:"item_id" binding:"required,uuid"`
	Tier   string `json:"tier" binding:"required,oneof=S A B C D F"`
}

type TierlistItemResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	ImageURL string `json:"image_url"`
	SortOrder int `json:"sort_order"`
}

type TierlistResponse struct {
	ID string `json:"id"`
	ShareCode string `json:"share_code"`
	Title string `json:"title"`
	Description string `json:"description"`
	ExpiresAt time.Time `json:"expires_at"`
	Creator UserReponse `json:"creator"`
	Items []TierlistItemResponse `json:"items"`
	HasSubmitted bool `json:"has_submitted"`
}

type CreateTierlistResponse struct {
	ShareCode string `json:"share_code"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TierlistResultResponse struct {
	Tierlist TierlistResponse `json:"tierlist"`
	TotalSubmissions int `json:"total_submissions"`
	Results []TierResult `json:"results"`
}

type TierResult struct {
	ItemID string `json:"item_id"`
	ItemName string `json:"item_name"`
	ImageURL string `json:"image_url"`
	Counts map[string]int `json:"counts"` // key is tier name, value is count of submissions for that tier
	Total int `json:"total"` // total number of submissions for this item
	TopTier string `json:"top_tier"` // the tier with the highest count, if tie then the higher tier wins (S > A > B > C > D > F)
	AverageScore float64 `json:"average_score"` // 1.0 - 6.0
	Rank int `json:"rank"` // rank of the item based on average score, starts at 1 
}