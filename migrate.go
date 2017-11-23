package main

import (
	"time"

	"github.com/go-gormigrate/gormigrate"
	"github.com/jinzhu/gorm"

	"github.com/zendesk/slack-poc/config"
	"github.com/zendesk/slack-poc/operation"
)

var versions = []*gormigrate.Migration{
	{
		ID: "20171123",
		Migrate: func(tx *gorm.DB) error {
			type Integration struct {
				SlackToken string `gorm:"primary_key;size:255;not null"`
				SlackWorkspace string `gorm:"size:255;not null"`
				ZendeskSubdomain int32 `gorm:"not null"`
				CreatedAt time.Time `gorm:"not null"`
			}

			if err := tx.AutoMigrate(
				&Integration{},
			).Error ; err != nil{
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("integration").Error
		},
	},
}

func main() {
	defaultOptions := &gormigrate.Options{
		TableName: "migration",
		IDColumnName: "id",
		UseTransaction: true,
	}

	cfg := config.GetConfig()
	operations.Migrate(cfg, versions, defaultOptions)
}
