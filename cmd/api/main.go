package main

import (
	"log"
	"time"

	"github.com/n-vh/twitch-renames/internal/config"
	"github.com/n-vh/twitch-renames/internal/database"
	"github.com/n-vh/twitch-renames/internal/graphql"
	"github.com/n-vh/twitch-renames/internal/helix"
	"github.com/n-vh/twitch-renames/internal/models"
	"github.com/n-vh/twitch-renames/internal/utils"
	"github.com/samber/lo"
)

type Service struct {
	workerIds []int
}

func main() {
	cfg := config.Get()
	database.Connect(&database.Config{
		Host:     cfg.PgHost,
		Port:     cfg.PgPort,
		User:     cfg.PgUser,
		Password: cfg.PgPassword,
		Database: cfg.PgDatabase,
		Schema:   cfg.PgSchema,
	})

	service := Service{}

	go service.init()
	utils.SetInterval(service.init, time.Hour)
}

func (service *Service) init() {
	count := database.Users.Count()

	for i := 0; i < count; i += config.MAX_PER_WORKER {
		workerId := int(float64(i / config.MAX_PER_WORKER))
		initialized := lo.Contains(service.workerIds, workerId)

		if !initialized {
			service.workerIds = append(service.workerIds, workerId)

			go func() {
				for {
					service.worker(workerId)
				}
			}()
		}
	}
}

func (service *Service) worker(workerId int) {
	initial := workerId * config.MAX_PER_WORKER
	max := initial + config.MAX_PER_WORKER

	worker := database.Workers.FindOneOrCreate(models.Worker{
		Cycles:   0,
		Offset:   initial,
		WorkerId: workerId,
	})

	for offset := worker.Offset; offset < max; offset += config.MAX_FETCH_DB {
		users := database.Users.FindAllPaginated(offset)

		if len(users) == 0 {
			break
		}

		for i := 0; i < len(users); i += config.MAX_CHUNK_SIZE {
			chunk := lo.Slice(users, i, i+config.MAX_CHUNK_SIZE)
			service.handleChunk(&chunk)
		}

		database.Workers.IncrementOffset(workerId)
	}

	database.Workers.ResetOne(models.Worker{
		Offset:   initial,
		WorkerId: workerId,
	})
}

func (service *Service) handleChunk(chunk *[]models.User) {
	userIds := lo.Map(*chunk, func(u models.User, i int) string {
		return u.UserId
	})

	data, ok := helix.GetUsersById(userIds)

	if !ok {
		return
	}

	lo.ForEach(*chunk, func(oldUser models.User, i int) {
		newUser, isRename := lo.Find(data, func(u helix.User) bool {
			return u.UserId == oldUser.UserId && u.Login != oldUser.Login
		})

		if !isRename {
			return
		}

		user, ok := graphql.GetUserLastUpdate(oldUser.UserId)

		if ok {
			isRecent := time.Since(user.UpdatedAt) < time.Hour*3
			date := lo.Ternary(isRecent, user.UpdatedAt, time.Now())

			ok := utils.InsertRename(oldUser, models.User{
				UserId:      newUser.UserId,
				Login:       newUser.Login,
				DisplayName: newUser.DisplayName,
			}, date)

			if ok {
				log.Printf("%s -> %s \n", oldUser.Login, newUser.Login)
			}
		}
	})
}
