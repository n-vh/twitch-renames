package main

import (
	"log"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/n-vh/twitch-renames/internal/config"
	"github.com/n-vh/twitch-renames/internal/database"
	"github.com/n-vh/twitch-renames/internal/graphql"
	"github.com/n-vh/twitch-renames/internal/helix"
	"github.com/n-vh/twitch-renames/internal/models"
	"github.com/n-vh/twitch-renames/internal/utils"
	"github.com/samber/lo"
)

type Connection struct {
	chat     *twitch.Client
	channels []string
}

type Service struct {
	Buffer      []models.User
	channels    []string
	connections []Connection
	mutex       sync.Mutex
	pause       bool
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

	service := Service{
		Buffer: []models.User{},
	}

	go service.init()
	utils.SetInterval(service.handleStreams, time.Minute*10)
}

func (service *Service) init() {
	service.channels = service.fetchChannels()

	for i := 0; i < len(service.channels); i += config.MAX_CHUNK_SIZE {
		go service.connection(lo.Slice(service.channels, i, i+config.MAX_CHUNK_SIZE))
		time.Sleep(time.Second * 2)
	}
}

func (service *Service) connection(channels []string) {
	chat := twitch.NewAnonymousClient()

	service.connections = append(service.connections, Connection{
		chat:     chat,
		channels: channels,
	})

	chat.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if !service.pause {
			service.handleMessage(message)
		}
	})

	chat.OnConnect(func() {
		log.Println("Connected")
		chat.Join(channels...)
	})

	chat.Connect()
}

func (service *Service) handleMessage(tags twitch.PrivateMessage) {
	userId, login, displayName := tags.User.ID, tags.User.Name, tags.User.DisplayName

	service.mutex.Lock()

	if len(service.Buffer) >= config.MAX_BUFFER_SIZE {
		service.pause = true

		go service.handleBuffer(service.Buffer)

		service.Buffer = []models.User{}
		service.pause = false
	} else {
		service.Buffer = append(service.Buffer, models.User{
			UserId:      userId,
			Login:       login,
			DisplayName: displayName,
		})
	}

	service.mutex.Unlock()
}

func (service *Service) handleBuffer(buffer []models.User) {
	values := lo.UniqBy(buffer, func(u models.User) string {
		return u.UserId
	})
	keys := lo.Map(values, func(u models.User, i int) string {
		return u.UserId
	})

	users := database.Users.FindIn(keys)

	newUsers := lo.Filter(values, func(newUser models.User, i int) bool {
		oldUser, found := lo.Find(users, func(oldUser models.User) bool {
			return oldUser.UserId == newUser.UserId
		})

		if found && oldUser.Login != newUser.Login {
			user, ok := graphql.GetUserLastUpdate(oldUser.UserId)

			if ok {
				isRecent := time.Since(user.UpdatedAt) < time.Hour*3
				date := lo.Ternary(isRecent, user.UpdatedAt, time.Now())

				ok := utils.InsertRename(oldUser, newUser, date)

				if ok {
					log.Printf("# %s -> %s \n", oldUser.Login, newUser.Login)
				}

			}
		} else if !found {
			return true
		}

		return false
	})

	if len(newUsers) > 0 {
		database.Users.InsertMany(&newUsers)
	}
}

func (service *Service) handleStreams() {
	channels := service.fetchChannels()
	parts, joins := lo.Difference(service.channels, channels)

	if len(parts) == 0 || len(joins) == 0 {
		return
	}

	service.channels = channels

	for _, client := range service.connections {
		parts := lo.Intersect(parts, client.channels)

		for _, part := range parts {
			var join string
			join, joins = joins[0], joins[1:]

			client.chat.Depart(part)
			client.channels = remove(client.channels, part)

			client.chat.Join(join)
			client.channels = append(client.channels, join)
		}
	}

	log.Printf("Updating %d channels\n", len(parts))
}

func (service *Service) fetchChannels() []string {
	streams := helix.GetStreamsPaginated(config.MAX_FETCH_STREAMS)

	users := lo.Map(streams, func(s helix.Stream, i int) models.User {
		return models.User{
			UserId:      s.UserId,
			Login:       s.Login,
			DisplayName: s.DisplayName,
		}
	})
	go service.handleBuffer(users)

	return lo.Map(streams, func(s helix.Stream, i int) string {
		return s.Login
	})
}

func remove[T comparable](slice []T, el T) []T {
	i := lo.IndexOf(slice, el)

	if i < 0 {
		return slice
	}

	return slice[:i+copy(slice[i:], slice[i+1:])]
}
