package models

type Tierlist struct {
	ID int `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatorID int `json:"creator_id"`
	Creator User `gorm:"foreignKey:CreatorID"`
	Tiers []Tiers `gorm:"foreignKey:TierlistID"`
	Items []Items `gorm:"foreignKey:TierlistID"`
	Version int `gorm:"default:1" json:"version"`
}

type Tiers struct {
	ID int `gorm:"primaryKey" json:"id"`
	TierlistID int `json:"tierlist_id"`
	Text string `json:"text"`
	Colour string `json:"colour"`
	Tierlist Tierlist `gorm:"foreignKey:TierlistID"`
	Items []Items `gorm:"foreignKey:TierID"`
}

type Items struct {
	ID int `gorm:"primaryKey" json:"id"`
	TierlistID int `json:"tierlist_id"`
	Text string `json:"text"`
	Image string `json:"image"`
	TierID int `json:"tier_id"`
	Tier Tiers `gorm:"foreignKey:TierID"`
	Tierlist Tierlist `gorm:"foreignKey:TierlistID"`
}