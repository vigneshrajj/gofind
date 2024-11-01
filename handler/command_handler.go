package handler

import (
	"github.com/vigneshrajj/gofind/models"
	"gorm.io/gorm"
)

func CreateCommand(db *gorm.DB, command models.Command) error {
	if err := db.Create(&command).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCommand(db *gorm.DB, alias string) error {
	if err := db.Delete(&models.Command{}, "alias=?", alias).Error; err != nil {
		return err
	}
	return nil
}

func ListCommands(db *gorm.DB) []models.Command {
	var commands []models.Command
	db.Find(&commands)
	return commands
}

func SearchCommand(db *gorm.DB, alias string, partialMatch bool) (models.Command, error) {
	var command models.Command
	if partialMatch {
		if err := db.Where("alias LIKE ?", alias+"%").Find(&command).Error; err != nil {
			return command, err
		}
	} else {
		if err := db.Where("alias=?", alias).Find(&command).Error; err != nil {
			return command, err
		}
	}
	return command, nil
}
