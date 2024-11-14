package database

import (
	"fmt"
	"tierlist/models"
)

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
		{Name: "Test 1", Email: "test1@test.com"},
		{Name: "Test 2", Email: "test2@test.com"},
		{Name: "Test 3", Email: "test3@test.com"},
	}

	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to seed user")
		}
	}

	
	fmt.Println("Users Seeded Successfully")
	return nil
}

func SeedTierlists() error {
	if DB == nil {
		return fmt.Errorf("database connection unavailable")
	}

	var tierlistResult []models.Tierlist
	result := DB.Find(&tierlistResult)

	if result.RowsAffected > 0 {
		return nil
	}

	tierlists := []models.Tierlist{
		{Name: "Tierlist1", Description: "This is tierlist 1", CreatorID: 1},
		{Name: "Tierlist2", Description: "This is tierlist 2", CreatorID: 2},
		{Name: "Tierlist3", Description: "This is tierlist 3", CreatorID: 3},
	}

	for _, tierlist := range tierlists {
		if err := DB.Create(&tierlist).Error; err != nil {
			return fmt.Errorf("failed to seed tierlists")
		}
	}

	fmt.Println("Tierlists Seeded Successfully")
	return nil
}

func SeedTiers() error {
	if DB == nil {
		return fmt.Errorf("database connection unavailable")
	}

	var tiersResult []models.Tiers
	result := DB.Find(&tiersResult)

	if result.RowsAffected > 0 {
		return nil
	}

	tiers := []models.Tiers{
		{TierlistID: 1, Text: "S", Colour: "f4ee1d"},
		{TierlistID: 1, Text: "A", Colour: "d81621"},
		{TierlistID: 1, Text: "B", Colour: "00bef1"},
		{TierlistID: 1, Text: "C", Colour: "37f100"},
		{TierlistID: 1, Text: "D", Colour: "ffffff"},
		{TierlistID: 1, Text: "F", Colour: "000000"},
	}

	for _, tier := range tiers {
		if err := DB.Create(&tier).Error; err != nil {
			return fmt.Errorf("failed to seed tiers")
		}
	}

	fmt.Println("Tiers Seeded Successfully")
	return nil
}

func SeedItems() error {
	if DB == nil {
		return fmt.Errorf("database connection unavailable")
	}

	var itemsResult []models.Items
	result := DB.Find(&itemsResult)

	if result.RowsAffected > 0 {
		return nil
	}

	items := []models.Items{
		{TierlistID: 1, Text: "Nacht der Untoten", Image: "https://static.wikia.nocookie.net/callofduty/images/2/2b/Nacht_Der_Untoten_Menu_Selection_WaW.png/revision/latest?cb=20161009103531", TierID: 4},
		{TierlistID: 1, Text: "Verrückt", Image: "https://static.wikia.nocookie.net/callofduty/images/a/a1/Verruckt_Menu_Selection_WaW.png/revision/latest?cb=20161009103542", TierID: 3},
		{TierlistID: 1, Text: "Shi No Numa", Image: "https://static.wikia.nocookie.net/callofduty/images/2/2f/Shi_No_Numa_Menu_Selection_WaW.png/revision/latest?cb=20161009103553", TierID: 2},
		{TierlistID: 1, Text: "Der Riese", Image: "https://static.wikia.nocookie.net/callofduty/images/8/86/Der_Riese_Menu_Selection_WaW.png/revision/latest?cb=20161009103603", TierID: 1},
		{TierlistID: 1, Text: "Moon", Image: "https://static.wikia.nocookie.net/callofduty/images/c/cc/Moon_Menu_Selection_BO.png/revision/latest?cb=20240710075602", TierID: 1},
	}

	for _, item := range items {
		if err := DB.Create(&item).Error; err != nil {
			return fmt.Errorf("failed to seed items")
		}
	}

	fmt.Println("Items Seeded Successfully")
	return nil
}