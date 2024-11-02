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
	if err := db.Delete(&models.Command{}, "alias=? AND is_default=?", alias, false).Error; err != nil {
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
		if err := db.Where("alias LIKE ?", alias+"%").Order("LENGTH(alias) ASC").Find(&command).Error; err != nil {
			return command, err
		}
	} else {
		if err := db.Where("alias=?", alias).Find(&command).Error; err != nil {
			return command, err
		}
	}
	return command, nil
}

func GetDefaultCommand(db *gorm.DB) (models.Command, error) {
	var command models.Command
	if err := db.Where("is_default=?", true).Find(&command).Error; err != nil {
		return command, err
	}
	return command, nil
}

func SetDefaultCommand(db *gorm.DB, alias string) error {
	var command models.Command
	var defaultCommand models.Command
	if err := db.Where("alias=?", alias).Find(&command).Error; err != nil {
		return err
	}
	if err := db.Where("is_default=?", true).Find(&defaultCommand).Error; err != nil {
		return err
	}
	command.IsDefault = true
	defaultCommand.IsDefault = false
	db.Save(&command)
	db.Save(&defaultCommand)
	return nil
}
