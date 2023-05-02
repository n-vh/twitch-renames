package utils

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/n-vh/twitch-renames/internal/database"
	"github.com/n-vh/twitch-renames/internal/models"
	"github.com/samber/lo"
)

func InsertRename(oldUser models.User, newUser models.User, date time.Time) bool {
	user, ok := database.Users.FindOne(oldUser.UserId)

	if ok && user.Login != newUser.Login {
		database.Users.UpdateOne(models.User{
			UserId:      newUser.UserId,
			Login:       newUser.Login,
			DisplayName: newUser.DisplayName,
		})

		database.Renames.InsertOne(models.Rename{
			UserId:         newUser.UserId,
			Login:          newUser.Login,
			DisplayName:    newUser.DisplayName,
			OldLogin:       oldUser.Login,
			OldDisplayName: oldUser.DisplayName,
			Date:           date,
		})

		return true
	}

	return false
}

func ParseAutoComplete(input string, renames []models.Rename) []fiber.Map {
	return lo.Map(renames, func(rename models.Rename, i int) fiber.Map {
		chooseRecent := strings.HasPrefix(rename.Login, input)

		login := lo.Ternary(chooseRecent, rename.Login, rename.OldLogin)
		displayName := lo.Ternary(chooseRecent, rename.DisplayName, rename.OldDisplayName)

		return fiber.Map{
			"login":       login,
			"displayName": displayName,
		}
	})
}

func ParseSearch(renames []models.Rename) []fiber.Map {
	return lo.Map(renames, func(rename models.Rename, i int) fiber.Map {
		return fiber.Map{
			"date": rename.Date.UTC().Format("2006-01-02T15:04:05.000Z"),
			"new": fiber.Map{
				"login":       rename.Login,
				"displayName": rename.DisplayName,
			},
			"old": fiber.Map{
				"login":       rename.OldLogin,
				"displayName": rename.OldDisplayName,
			},
			"userId": rename.UserId,
		}
	})
}

// func ParseJsonBody[R any](body io.ReadCloser) (R, bool) {
// 	var result R

// 	bytes, _ := ioutil.ReadAll(body)
// 	err := json.Unmarshal(bytes, &result)

// 	if err != nil {
// 		log.Println("PARSE ERROR", result, err)
// 		return result, false
// 	}

// 	return result, true
// }

func ParseJsonBody[T any](b io.ReadCloser) (r T, ok bool) {
	err := json.NewDecoder(b).Decode(&r)
	if err != nil {
		log.Println("PARSE ERROR", r, err)
	}
	return r, err == nil
}

func SetInterval(callback func(), i time.Duration) {
	for {
		time.Sleep(i)
		go callback()
	}
}
