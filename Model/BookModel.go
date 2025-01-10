package Model

type BookModel struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title" gorm:"not null"`
	Author      string `json:"author" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	Available   bool   `json:"available" gorm:"not null"`
}
