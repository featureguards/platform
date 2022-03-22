package models

type User struct {
	Model
	OryID    string    `gorm:"uniqueIndex"`
	Projects []Project `gorm:"foreignKey:OwnerID"`
}
