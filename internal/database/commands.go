package database

import (
	"errors"

	"github.com/vigneshrajj/gofind/models"
	"gorm.io/gorm"
)

func CreateCommand(db *gorm.DB, command models.Command) error {
	if err := db.Create(&command).Error; err != nil {
		return err
	}
	return nil
}

func FirstOrCreateCommand(db *gorm.DB, command models.Command) {
	db.FirstOrCreate(&command)
}

func DeleteCommand(db *gorm.DB, alias string) error {
	if rowsAffected := db.Delete(&models.Command{}, "alias=? AND is_default=?", alias, false).RowsAffected; rowsAffected == 0 {
		return errors.New("Command not found")
	}
	return nil
}

func ListCommands(db *gorm.DB) []models.Command {
	var commands []models.Command
	db.Find(&commands)
	return commands
}

func FilteredListCommands(db *gorm.DB, query string, pageSize int, offset int) (*[]models.Command, error) {
    var commands []models.Command

    switch {
    case pageSize > 100:
        pageSize = 100
    case pageSize <= 0:
        pageSize = 10
    }

    if query != "" {
        db = db.Where("alias LIKE ? OR query LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
    }

    // db = db.Offset(offset).Limit(pageSize).Order("type ASC")

    if err := db.Find(&commands).Error; err != nil {
        return nil, err
    }

    return &commands, nil
}

func SearchCommand(db *gorm.DB, alias string, partialMatch bool) models.Command {
	var command models.Command
	if partialMatch {
		db.Where("alias LIKE ?", alias+"%").Order("LENGTH(alias) ASC").Find(&command)
	} else {
		db.Where("alias=?", alias).Find(&command)
	}
	return command
}

func GetDefaultCommand(db *gorm.DB) models.Command {
	var command models.Command
	db.Where("is_default=?", true).Find(&command)
	return command
}

func SetDefaultCommand(db *gorm.DB, alias string) error {
	var command models.Command
	var defaultCommand models.Command
	if db.Where("alias=?", alias).Find(&command); command == (models.Command{}) {
		return errors.New("Command not found")
	}
	if db.Where("is_default=?", true).Find(&defaultCommand); defaultCommand == (models.Command{}) {
		return errors.New("Default Command not found")
	}
	command.IsDefault = true
	defaultCommand.IsDefault = false
	db.Save(&command)
	db.Save(&defaultCommand)
	return nil
}
