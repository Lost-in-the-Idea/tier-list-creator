package models

type Tierlist struct {
	ID int `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatorID int `json:"creator_id"`
	Creator User `gorm:"foreignKey:CreatorID"`
	Tiers []Tier `gorm:"foreignKey:TierlistID"`
	Items []Item `gorm:"foreignKey:TierlistID"`
	Version int `gorm:"default:1" json:"version"`
}

type Tier struct {
	ID int `gorm:"primaryKey" json:"id"`
	TierlistID int `json:"tierlist_id"`
	Text string `json:"text"`
	Colour string `json:"colour"`
	Tierlist Tierlist `gorm:"foreignKey:TierlistID"`
	// Items []Item `gorm:"foreignKey:TierID"` // This is not needed for now
}

type Item struct {
	ID int `gorm:"primaryKey" json:"id"`
	TierlistID int `json:"tierlist_id"`
	Text string `json:"text"`
	Image string `json:"image"`
	// TierID int `json:"tier_id"` // This is not needed for now
	TierText string `gorm:"default:U" json:"tier_text"`
	// Tier Tier `gorm:"foreignKey:TierID"` // This is not needed for now
	Tierlist Tierlist `gorm:"foreignKey:TierlistID"`
}

// POST & PUT /tierlist structs

type TierRequest struct {
	Name string `json:"text" binding:"required"`
	Colour string `json:"colour" binding:"required"`
}

type ItemRequest struct {
	Text string `json:"text" binding:"required"`
	Image string `json:"image" binding:"required"`
	Tier string `json:"tier" binding:"required"`
}	

type TierlistRequest struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Tiers []TierRequest `json:"tiers" binding:"required"`
	Items []ItemRequest `json:"items" binding:"required"`
	Creator int `json:"creator_id" binding:"required"` // placeholder for now
}

type UpdateTierlistRequest struct {
	Tiers []TierRequest `json:"tiers" binding:"required"`
	Items []ItemRequest `json:"items" binding:"required"`
}