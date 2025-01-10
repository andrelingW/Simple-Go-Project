package Model

type UserModel struct {
	UserId   int    `json:"userId" gorm:"primaryKey;autoIncrement"`
	UserName string `json:"userName" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"unique;not null"`
}
