package diary_dao

import (
	"context"
	"fmt"
	"github.com/majolo/web-app-starter/services/diary_dao/migration"
	"time"

	database "github.com/majolo/web-app-starter/database"
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	UserId    string
	Text      string
	CreatedAt time.Time
}

type DiaryDAO struct {
	db *gorm.DB
}

func NewDiaryDAO(db *gorm.DB) (DiaryDAO, error) {
	err := database.Migrate(db, migration.MigrationsDiary)
	if err != nil {
		return DiaryDAO{}, err
	}
	return DiaryDAO{
		db: db,
	}, nil
}

func (d *DiaryDAO) CreateDiaryEntry(ctx context.Context, entry *Entry) (uint, error) {
	if entry == nil {
		return 0, fmt.Errorf("no entry provided")
	}
	entry.CreatedAt = time.Now()
	tx := d.db.Begin().WithContext(ctx)
	err := tx.Create(entry).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = tx.Commit().Error
	if err != nil {
		return 0, err
	}
	return entry.ID, nil
}

func (d *DiaryDAO) ListDiaryEntries(ctx context.Context, userId string) ([]*Entry, error) {
	var entries []*Entry
	err := d.db.WithContext(ctx).
		Where(&Entry{UserId: userId}).
		Find(&entries).Error
	if err != nil {
		return nil, err
	}
	return entries, nil
}
