package db

import (
	"fmt"
	"log"

	model "backend/internal/models"
)

func GetIdea(baseIdeaID uint64) (*model.BaseIdea, error) {
	var baseIdea model.BaseIdea

	// 1. BaseIdeaに紐づくCardsを取得
	if err := db.Preload("Cards").First(&baseIdea, baseIdeaID).Error; err != nil {
		log.Printf("Error loading BaseIdea with ID %d: %v", baseIdeaID, err)
		return nil, fmt.Errorf("failed to load BaseIdea with ID %d: %w", baseIdeaID, err)
	}

	// 2. 各Cardに紐づくIdeaを再帰的に取得
	for i := range baseIdea.Cards {
		if err := preloadIdeas(&baseIdea.Cards[i]); err != nil {
			return nil, err
		}
	}

	return &baseIdea, nil
}

func preloadIdeas(card *model.Card) error {
	// Cardに紐づくIdeaを取得
	if err := db.Preload("Ideas.Cards").Find(&card).Error; err != nil {
		return fmt.Errorf("failed to preload Ideas for Card with ID %d: %v", card.ID, err)
	}

	// 各Ideaに紐づくCardを再帰的に取得
	for i := range card.Ideas {
		if err := preloadCards(&card.Ideas[i]); err != nil {
			return err
		}
	}

	return nil
}

func preloadCards(idea *model.Idea) error {
	// Ideaに紐づくCardを取得
	if err := db.Preload("Cards.Ideas").Find(&idea).Error; err != nil {
		return fmt.Errorf("failed to preload Ideas for Card with ID %d: %v", idea.ID, err)
	}

	// 各Cardに紐づくIdeaを再帰的に取得
	for i := range idea.Cards {
		if err := preloadIdeas(&idea.Cards[i]); err != nil {
			return err
		}
	}

	return nil
}

func ChangeCardStatus(cardID uint64, cardStatus uint) error {
	var card model.Card
	if err := db.Model(&card).Where("id = ?", cardID).Update("status_flag", cardStatus).Error; err != nil {
		return err
	}

	return nil
}

func CreateIdea(ideaTitles []string, cardID uint) error {

	type IdeaInfo struct {
		ID    uint   `gorm:"primary_key;auto_increment" json:"id"`
		Title string `json:"title"`
	}

	var ideaRecords []IdeaInfo

	for _, ideaTitle := range ideaTitles {
		ideaRecords = append(ideaRecords, IdeaInfo{Title: ideaTitle})
	}
	if err := db.Table("ideas").Create(&ideaRecords).Error; err != nil {
		return err
	}

	type CardIdeas struct {
		CardID uint `json:"card_id"`
		IdeaID uint `json:"idea_id"`
	}

	var cardIdeas []CardIdeas

	for _, ideaInfo := range ideaRecords {
		cardIdeas = append(cardIdeas, CardIdeas{
			CardID: cardID,
			IdeaID: ideaInfo.ID,
		})
	}
	if err := db.Create(&cardIdeas).Error; err != nil {
		return err
	}
	return nil
}

func CreateCard(cardTitle string, baseIdeaID uint, ideaID uint) (uint, error) {
	var cardInfo model.Card
	cardInfo.BaseIdeaID = baseIdeaID
	cardInfo.Title = cardTitle
	cardInfo.StatusFlag = 1

	if err := db.Create(&cardInfo).Error; err != nil {
		return 0, err
	}

	type IdeaCard struct {
		IdeaID uint `json:"idea_id"`
		CardID uint `json:"card_id"`
	}

	var ideaCard IdeaCard
	if baseIdeaID == 0 {

		ideaCard.IdeaID = ideaID
		ideaCard.CardID = cardInfo.ID
		if err := db.Create(&ideaCard).Error; err != nil {
			return 0, err
		}
	}

	return cardInfo.ID, nil
}

func GetCard(cardID uint) (*model.Card, error) {
	var card model.Card
	if err := db.Preload("Ideas").Where("id = ?", cardID).First(&card).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func DeleteCard(cardID uint) error {
	var card model.Card
	if err := db.First(&card, cardID).Error; err != nil {
		return err
	}

	var ideas []model.Idea
	db.Model(&card).Association("Ideas").Find(&ideas)

	for _, idea := range ideas {
		deleteIdeaRecursively(&idea)
	}

	if err := db.Model(&card).Association("Ideas").Clear(); err != nil {
		return err
	}

	if err := db.Table("idea_cards").Where("card_id = ?", cardID).Delete(nil).Error; err != nil {
		return err
	}

	if err := db.Delete(&card).Error; err != nil {
		return err
	}

	return nil
}
