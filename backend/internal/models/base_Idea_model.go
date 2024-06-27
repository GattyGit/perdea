package models

type BaseIdea struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	UserID uint   `json:"user_id"`
	Cards  []Card `gorm:"foreignKey:BaseIdeaID"`
}
