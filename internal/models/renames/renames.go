package renames

import (
	"github.com/n-vh/twitch-renames/internal/models"
	"gorm.io/gorm"
)

type Model struct {
	db *gorm.DB
}

func InitModel(db *gorm.DB) *Model {
	db.AutoMigrate(&models.Rename{})
	db.Exec("CREATE INDEX IF NOT EXISTS login_like_idx ON \"renames\" (login COLLATE \"C\")")
	db.Exec("CREATE INDEX IF NOT EXISTS old_login_like_idx ON \"renames\" (old_login COLLATE \"C\")")
	return &Model{db: db}
}

func (ctx *Model) InsertOne(u models.Rename) error {
	return ctx.db.Create(&u).Error
}
