package helix

import (
	"log"
	"net/http"
	"strings"

	"github.com/n-vh/twitch-renames/internal/utils"
	"github.com/samber/lo"
)

func GetUsersById(userIds []string) ([]User, bool) {
	params := strings.Join(userIds, "&id=")

	url := ENDPOINT + "/users?id=" + params
	res, err := http.Get(url)

	if err != nil {
		log.Panic(err)
	}

	data, ok := utils.ParseJsonBody[HelixUserData](res.Body)

	if !ok {
		return nil, false
	}

	users := lo.Map(data.Data, func(u HelixUser, i int) User {
		return User(u)
	})

	return users, true
}
