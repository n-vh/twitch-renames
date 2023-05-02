package workers

import (
	"github.com/n-vh/twitch-renames/internal/config"
	"github.com/n-vh/twitch-renames/internal/models"
	"gorm.io/gorm"
)

type Model struct {
	db *gorm.DB
}

func InitModel(db *gorm.DB) *Model {
	db.AutoMigrate(&models.Worker{})
	return &Model{db: db}
}

func (ctx *Model) FindOne(workerId int) (models.Worker, bool) {
	var worker models.Worker

	err := ctx.db.First(&worker, "worker_id = ?", workerId).Error

	return worker, err == nil
}

func (ctx *Model) FindOneOrCreate(w models.Worker) models.Worker {
	worker, ok := ctx.FindOne(w.WorkerId)

	if !ok {
		worker = models.Worker(w)
		ctx.db.Create(&worker)
	}

	return worker
}

func (ctx *Model) IncrementOffset(workerId int) error {
	worker, _ := ctx.FindOne(workerId)

	return ctx.db.Where("worker_id = ?", workerId).Select("offset").Updates(models.Worker{
		Offset: worker.Offset + config.MAX_FETCH_DB,
	}).Error
}

func (ctx *Model) ResetOne(w models.Worker) error {
	return ctx.db.Exec("UPDATE workers SET \"cycles\" = \"cycles\" + 1, \"offset\" = ? WHERE worker_id = ?", w.Offset, w.WorkerId).Error
}
