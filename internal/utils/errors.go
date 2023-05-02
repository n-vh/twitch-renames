package utils

import "github.com/gofiber/fiber/v2"

func SendErrorBadRequest(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
		"status":  fiber.ErrBadRequest.Code,
		"message": fiber.ErrBadRequest.Message,
		"data":    []fiber.Map{},
	})
}

func SendErrorNotFound(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
		"status":  fiber.ErrNotFound.Code,
		"message": fiber.ErrNotFound.Message,
		"data":    []fiber.Map{},
	})
}
