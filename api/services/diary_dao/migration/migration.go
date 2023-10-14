package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"time"
)

var MigrationsDiary = []*gormigrate.Migration{
	{
		ID: "0001_initial_diary",
		Migrate: func(db *gorm.DB) error {
			// Copy the struct to avoid side effects and preserves original changes.
			type Entry struct {
				gorm.Model
				UserId    string
				Text      string
				CreatedAt time.Time
			}

			if err := db.AutoMigrate(&Entry{}); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(db *gorm.DB) error {
			type Entry struct{}
			if err := db.Migrator().DropTable(Entry{}); err != nil {
				return err
			}
			return nil
		},
	},
}
