package models

type Idea struct {
	ID    uint   `json:"id" gorm:"primary_key;auto_increment"`
	Title string `json:"title"`
	Cards []Card `gorm:"many2many:idea_cards;"`
}
