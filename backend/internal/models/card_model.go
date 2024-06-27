package models

type Card struct {
	ID         uint   `json:"id"`
	BaseIdeaID uint   `json:"base_idea_id"`
	Title      string `json:"title"`
	StatusFlag uint   `json:"status_flag"`
	Ideas      []Idea `gorm:"many2many:card_ideas;"`
}
