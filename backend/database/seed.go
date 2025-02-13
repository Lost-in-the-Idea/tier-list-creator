package database

import (
	"fmt"
	"tierlist/models"
	"time"
)

func Seed() {
	SeedUsers()
	SeedTierlist()
}

func SeedTierlist() error {
	if DB == nil {
		return fmt.Errorf("database connection unavailable")
	}

	var tierlistResult []models.Tierlist
	result := DB.Find(&tierlistResult)

	if result.RowsAffected > 0 {
		return nil
	}

	tierlists := []models.Tierlist{
		{Name: "Tierlist 1",
		Description:  "This is the first Tierlist",
		CreatorID: 1,
		Tiers: []models.Tier{
			{Text: "S", Colour: "f4ee1d"},
			{Text: "A", Colour: "d81621"},
			{Text: "B", Colour: "00bef1"},
			{Text: "C", Colour: "37f100"},
			{Text: "D", Colour: "ffffff"},
			{Text: "F", Colour: "000000"},
		},
		Items: []models.Item{
		{Text: "Nacht der Untoten", Image: "https://static.wikia.nocookie.net/callofduty/images/2/2b/Nacht_Der_Untoten_Menu_Selection_WaW.png/revision/latest?cb=20161009103531", TierText: "S"},
		{Text: "Verrückt", Image: "https://static.wikia.nocookie.net/callofduty/images/a/a1/Verruckt_Menu_Selection_WaW.png/revision/latest?cb=20161009103542", TierText: "A"},
		{Text: "Shi No Numa", Image: "https://static.wikia.nocookie.net/callofduty/images/2/2f/Shi_No_Numa_Menu_Selection_WaW.png/revision/latest?cb=20161009103553", TierText: "B"},
		{Text: "Der Riese", Image: "https://static.wikia.nocookie.net/callofduty/images/8/86/Der_Riese_Menu_Selection_WaW.png/revision/latest?cb=20161009103603", TierText: "C"},
		{Text: "Moon", Image: "https://static.wikia.nocookie.net/callofduty/images/c/cc/Moon_Menu_Selection_BO.png/revision/latest?cb=20240710075602", TierText: "U"},
		},
		Version: 1,},
		{Name: "Tierlist 2",
		Description:  "This is the second Tierlist",
		CreatorID: 2,
		Tiers: []models.Tier{
			{Text: "S", Colour: "f4ee1d"},
			{Text: "A", Colour: "d81621"},
			{Text: "B", Colour: "00bef1"},
			{Text: "C", Colour: "37f100"},
			{Text: "D", Colour: "ffffff"},
			{Text: "F", Colour: "000000"},
		},
		Items: []models.Item{
			{Text: "Green Hill Zone", Image: "https://example.com/green_hill_zone.png", TierText: "S"},
			{Text: "Chemical Plant Zone", Image: "https://example.com/chemical_plant_zone.png", TierText: "A"},
			{Text: "Casino Night Zone", Image: "https://example.com/casino_night_zone.png", TierText: "B"},
			{Text: "Ice Cap Zone", Image: "https://example.com/ice_cap_zone.png", TierText: "C"},
			{Text: "Sky Sanctuary Zone", Image: "https://example.com/sky_sanctuary_zone.png", TierText: "U"},
		},
		Version: 1,},
		{Name: "Tierlist 3",
		Description:  "This is the third Tierlist",
		CreatorID: 3,
		Tiers: []models.Tier{
			{Text: "S", Colour: "f4ee1d"},
			{Text: "A", Colour: "d81621"},
			{Text: "B", Colour: "00bef1"},
			{Text: "C", Colour: "37f100"},
			{Text: "D", Colour: "ffffff"},
			{Text: "F", Colour: "000000"},
		},
		Items: []models.Item{
			{Text: "Mario", Image: "https://example.com/mario.png", TierText: "S"},
			{Text: "Luigi", Image: "https://example.com/luigi.png", TierText: "A"},
			{Text: "Peach", Image: "https://example.com/peach.png", TierText: "B"},
			{Text: "Bowser", Image: "https://example.com/bowser.png", TierText: "C"},
			{Text: "Yoshi", Image: "https://example.com/yoshi.png", TierText: "U"},
		},
		Version: 1,
		},
	}
	for _, tierlist := range tierlists {
	if err := DB.Create(&tierlist).Error; err != nil {
		return fmt.Errorf("failed to seed tierlist")
	}
}

	fmt.Println("Tierlists Seeded Successfully")
	return nil
}

func SeedUsers() error {
	if DB == nil {
		return fmt.Errorf("database connection unavailable")
	}

	var userResult []models.User
	result := DB.Find(&userResult)

	if result.RowsAffected > 0 {
		return nil
	}

	users := []models.User{
        {
            DiscordID: "123456789",
            Username:  "Test User 1",
            Avatar:    "default1",
            LastLogin: time.Now(),
        },
        {
            DiscordID: "987654321",
            Username:  "Test User 2",
            Avatar:    "default2",
            LastLogin: time.Now(),
        },
        {
            DiscordID: "456789123",
            Username:  "Test User 3",
            Avatar:    "default3",
            LastLogin: time.Now(),
        },
	}

	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to seed user")
		}
	}

	
	fmt.Println("Users Seeded Successfully")
	return nil
}