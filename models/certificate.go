package models

type Certificate struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"unique;not null"`
	SignRequest string
	Certificate string
	Token       string `gorm:"unique;not null"`
}