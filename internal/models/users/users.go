package users

import (
	"github.com/n-vh/twitch-renames/internal/config"
	"github.com/n-vh/twitch-renames/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model struct {
	db *gorm.DB
}

func InitModel(db *gorm.DB) *Model {
	db.AutoMigrate(&models.User{})
	return &Model{db: db}
}

func (ctx *Model) Count() int {
	var count int64

	ctx.db.Model(models.User{}).Count(&count)

	return int(count)
}

func (ctx *Model) FindAllPaginated(offset int) []models.User {
	users := []models.User{}

	ctx.db.Where("id >= ?", offset).Order("id ASC").Limit(config.MAX_FETCH_DB).Find(&users)

	return users
}

func (ctx *Model) FindOne(userId string) (models.User, bool) {
	var user models.User

	err := ctx.db.First(&user, "user_id = ?", userId).Error

	return user, err == nil
}

func (ctx *Model) FindIn(slice []string) []models.User {
	users := []models.User{}

	ctx.db.Find(&users, "user_id IN ?", slice)

	return users
}

func (ctx *Model) InsertMany(users *[]models.User) {
	ctx.db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&users)
}

func (ctx *Model) UpdateOne(u models.User) error {
	return ctx.db.Where("user_id = ?", u.UserId).Updates(&u).Error
}
