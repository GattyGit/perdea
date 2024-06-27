package db

import (
	model "backend/internal/models"
)

func SaveBaseIdea(baseIdea *model.BaseIdea) (uint, error) {
	var baseIdeaID uint

	if err := db.Create(&baseIdea).Error; err != nil {
		return 0, err
	}

	baseIdeaID = baseIdea.ID

	return baseIdeaID, nil
}

func GetBaseIdea(user_id int) (*[]model.BaseIdea, error) {
	var baseIdeas []model.BaseIdea

	err := db.Where("user_id = ?", user_id).Find(&baseIdeas).Error

	if err != nil {
		return nil, err
	}
	return &baseIdeas, nil
}

func DeleteBaseIdeaRecursively(baseIdeaID uint) error {
	var baseIdea model.BaseIdea
	if err := db.First(&baseIdea, baseIdeaID).Error; err != nil {
		return err
	}

	var cards []model.Card
	db.Model(&baseIdea).Association("Cards").Find(&cards)

	for _, card := range cards {
		deleteCardRecursively(&card)
	}

	if err := db.Delete(&baseIdea).Error; err != nil {
		return err
	}
	return nil
}

func deleteCardRecursively(card *model.Card) {
	var ideas []model.Idea
	db.Model(card).Association("Ideas").Find(&ideas)

	for _, idea := range ideas {
		deleteIdeaRecursively(&idea)
	}

	db.Model(card).Association("Ideas").Clear()

	db.Delete(card)
}

func deleteIdeaRecursively(idea *model.Idea) {
	var cards []model.Card
	db.Model(idea).Association("Cards").Find(&cards)

	for _, card := range cards {
		deleteCardRecursively(&card)
	}

	db.Model(idea).Association("Cards").Clear()

	db.Delete(idea)
}
