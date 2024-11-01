package handler

import (
	"log"

	"github.com/vigneshrajj/gofind/models"
	"gorm.io/gorm"
)

func CreateCommand(db *gorm.DB, command models.Command) error {
	if err := db.Create(&command).Error; err != nil {
		log.Fatalf("Failed to create the command: %v", err)
	}
	return nil
}

func DeleteCommand(db *gorm.DB, id uint) error {
	if err := db.Delete(&models.Command{}, id).Error; err != nil {
		log.Fatalf("Failed to delete the command: %v", err)
	}
	return nil
}

func ListCommands(db *gorm.DB) []models.Command {
	var commands []models.Command
	db.Find(&commands)
	return commands
}

func SearchCommand(db *gorm.DB, alias string) []models.Command {
	var commands []models.Command
	db.Where("alias LIKE ?", alias+"%").Find(&commands)
	return commands
}
