package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/n-vh/twitch-renames/internal/database"
	"github.com/n-vh/twitch-renames/internal/utils"
)

func AutoComplete(ctx *fiber.Ctx) error {
	username, ok := utils.SanitizeUsername(ctx.Params("username"))

	if !ok {
		return utils.SendErrorBadRequest(ctx)
	}

	renames := database.Renames.AutoComplete(username)

	if len(renames) == 0 {
		return utils.SendErrorNotFound(ctx)
	}

	return ctx.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   utils.ParseAutoComplete(username, renames),
	})
}
