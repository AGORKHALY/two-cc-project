package models

import "gorm.io/gorm"

type Cars struct {
	ID      uint    `gorm:"primaryKey; autoIncrement" json:"id"`
	Company *string `json:"company"`
	Model   *string `json:"model"`
	Color   *string `json:"color"`
}

func MigrateCars(db *gorm.DB) error {
	err := db.AutoMigrate(&Cars{})

	return err
}
