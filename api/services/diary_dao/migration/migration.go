package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var MigrationsDiary = []*gormigrate.Migration{
	{
		ID: "0001_initial_diary",
		Migrate: func(db *gorm.DB) error {
			// Copy the struct to avoid side effects and preserves original changes.
			type Version struct {
				gorm.Model
				ModelMetadataID uint
				Version         string
				VersionUriPath  string // Path to the model folder within specified remote.
				RCloneRemote    string // Remote name from configured remotes.
			}
			if err := db.AutoMigrate(&Version{}); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(db *gorm.DB) error {
			type Version struct{}
			if err := db.Migrator().DropTable(Version{}); err != nil {
				return err
			}
			return nil
		},
	},
}
