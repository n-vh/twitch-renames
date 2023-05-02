package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/n-vh/twitch-renames/internal/database"
	"github.com/n-vh/twitch-renames/internal/utils"
)

func User(ctx *fiber.Ctx) error {
	id, ok := utils.SanitizeId(ctx.Params("id"))

	if !ok {
		return utils.SendErrorBadRequest(ctx)
	}

	user, ok := database.Users.FindOne(id)

	if !ok {
		return utils.SendErrorNotFound(ctx)
	}

	return ctx.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data": fiber.Map{
			"userId":     user.UserId,
			"login":      user.Login,
			"displaName": user.DisplayName,
		},
	})
}
