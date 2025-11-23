package model

type Team struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"size:255;not null"`
	Members []User `gorm:"many2many:user_team"`
}
