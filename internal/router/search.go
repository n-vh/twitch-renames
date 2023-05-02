package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/n-vh/twitch-renames/internal/database"
	"github.com/n-vh/twitch-renames/internal/utils"
)

func Search(ctx *fiber.Ctx) error {
	username, ok := utils.SanitizeUsername(ctx.Params("username"))

	if !ok {
		return utils.SendErrorBadRequest(ctx)
	}

	userIds := database.Renames.PreliminarySearch(username)

	if len(userIds) == 0 {
		return utils.SendErrorNotFound(ctx)
	}

	renames := database.Renames.Search(userIds)

	return ctx.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   utils.ParseSearch(renames),
	})
}
