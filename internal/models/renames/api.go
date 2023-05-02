package renames

import (
	"fmt"

	"github.com/n-vh/twitch-renames/internal/models"
	"github.com/samber/lo"
)

func (ctx *Model) AutoComplete(username string) []models.Rename {
	renames := []models.Rename{}

	filter := fmt.Sprintf("LIKE '%s'", username+"%")
	where := fmt.Sprintf("login %s OR old_login %s", filter, filter)

	ctx.db.Where(where).Order("id DESC").Limit(5).Find(&renames)

	return renames
}

func (ctx *Model) PreliminarySearch(username string) []string {
	renames := []models.Rename{}

	filter := fmt.Sprintf("LIKE '%s'", username+"%")
	where := fmt.Sprintf("login %s OR old_login %s", filter, filter)

	ctx.db.Select("user_id").Where(where).Order("id ASC").Find(&renames)

	return lo.Map(renames, func(r models.Rename, i int) string {
		return r.UserId
	})
}

func (ctx *Model) Search(userIds []string) []models.Rename {
	renames := []models.Rename{}

	ctx.db.Limit(100).Order("id DESC").Find(&renames, "user_id IN ?", userIds)

	return renames
}
