package services

import (
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"tierlist/database/models"
	"tierlist/dto"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
	ErrConflict  = errors.New("already submitted")
)

type TierlistService struct {
	db *gorm.DB
}

func NewTierlistService(db *gorm.DB) *TierlistService {
	return &TierlistService{db: db}
}

func (s *TierlistService) GetByID(id string, userID *uuid.UUID) (*dto.TierlistResponse, error) {
	var tierlist models.Tierlist
	if err := s.db.Where("id = ?", id).Preload("Creator").Preload("TierlistItems").First(&tierlist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	hasSubmitted := false
	if userID != nil {
		var submission models.Submissions
		if err := s.db.Where("tierlist_id = ? AND user_id = ?", tierlist.ID, *userID).First(&submission).Error; err == nil {
			hasSubmitted = true
		}
	}

	return &dto.TierlistResponse{
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
		Items:        mapTierlistItems(tierlist.TierlistItems),
		HasSubmitted: hasSubmitted,
	}, nil
}

func (s *TierlistService) Create(req dto.CreateTierlistRequest, creatorID uuid.UUID) (*dto.CreateTierlistResponse, error) {
	var result dto.CreateTierlistResponse
	err := s.db.Transaction(func(tx *gorm.DB) error {
		tierlist := models.Tierlist{
			Title:       req.Title,
			Description: req.Description,
			CreatorID:   creatorID,
			ShareCode:   uuid.New().String(),
			ExpiryTime:  req.ExpiryTime,
		}
		if err := tx.Create(&tierlist).Error; err != nil {
			return err
		}
		for _, item := range req.Items {
			newItem := models.TierlistItem{
				TierlistID: tierlist.ID,
				Name:       item.Name,
				ImageURL:   item.ImageURL,
				SortOrder:  item.SortOrder,
			}
			if err := tx.Create(&newItem).Error; err != nil {
				return err
			}
		}
		result = dto.CreateTierlistResponse{
			ShareCode: tierlist.ShareCode,
			ExpiresAt: tierlist.ExpiryTime,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *TierlistService) Submit(id string, userID uuid.UUID, req dto.SubmitRankingRequest) error {
	var tierlist models.Tierlist
	if err := s.db.Where("id = ?", id).Preload("TierlistItems").First(&tierlist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	var existing models.Submissions
	if err := s.db.Where("tierlist_id = ? AND user_id = ?", tierlist.ID, userID).First(&existing).Error; err == nil {
		return ErrConflict
	}

	validItemIDs := make(map[uuid.UUID]bool, len(tierlist.TierlistItems))
	for _, item := range tierlist.TierlistItems {
		validItemIDs[item.ID] = true
	}

	type parsedRanking struct {
		itemID uuid.UUID
		tier   string
	}
	rankings := make([]parsedRanking, 0, len(req.Rankings))
	for _, r := range req.Rankings {
		itemID, err := uuid.Parse(r.ItemID)
		if err != nil || !validItemIDs[itemID] {
			return fmt.Errorf("%w: invalid item ID %s", ErrNotFound, r.ItemID)
		}
		rankings = append(rankings, parsedRanking{itemID, r.Tier})
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		submission := models.Submissions{TierlistID: tierlist.ID, UserID: userID}
		if err := tx.Create(&submission).Error; err != nil {
			return err
		}
		for _, r := range rankings {
			rec := models.SubmissionRankings{
				SubmissionID: submission.ID,
				ItemID:       r.itemID,
				Tier:         r.tier,
			}
			if err := tx.Create(&rec).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *TierlistService) GetResults(id string) (*dto.TierlistResultResponse, error) {
	var tierlist models.Tierlist
	if err := s.db.Where("id = ?", id).Preload("Creator").Preload("TierlistItems").Preload("Submissions.Rankings").First(&tierlist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
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
		st := statsMap[item.ID]

		totalScore := 0
		for tier, count := range st.counts {
			totalScore += tierScores[tier] * count
		}
		var avgScore float64
		if st.total > 0 {
			avgScore = float64(totalScore) / float64(st.total)
		}

		topTier := ""
		topCount := 0
		for _, tier := range tierPrecedence {
			if st.counts[tier] > topCount {
				topCount = st.counts[tier]
				topTier = tier
			}
		}

		results = append(results, dto.TierResult{
			ItemID:       item.ID.String(),
			ItemName:     item.Name,
			ImageURL:     item.ImageURL,
			Counts:       st.counts,
			Total:        st.total,
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

	return &dto.TierlistResultResponse{
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
			Items:        mapTierlistItems(tierlist.TierlistItems),
			HasSubmitted: false,
		},
		TotalSubmissions: len(tierlist.Submissions),
		Results:          results,
	}, nil
}

func (s *TierlistService) Delete(id string, requestingUserID uuid.UUID) error {
	var tierlist models.Tierlist
	if err := s.db.Where("id = ?", id).First(&tierlist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if tierlist.CreatorID != requestingUserID {
		return ErrForbidden
	}
	return s.db.Delete(&tierlist).Error
}

func mapTierlistItem(item models.TierlistItem) dto.TierlistItemResponse {
	return dto.TierlistItemResponse{
		ID:        item.ID.String(),
		Name:      item.Name,
		ImageURL:  item.ImageURL,
		SortOrder: item.SortOrder,
	}
}

func mapTierlistItems(items []models.TierlistItem) []dto.TierlistItemResponse {
	out := make([]dto.TierlistItemResponse, len(items))
	for i, item := range items {
		out[i] = mapTierlistItem(item)
	}
	return out
}
