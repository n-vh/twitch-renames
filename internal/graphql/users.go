package graphql

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/n-vh/twitch-renames/internal/utils"
)

func GetUserLastUpdate(userId string) (LastUpdatedUser, bool) {
	json := []byte(`{"operationName":"UserLastUpdate","variables":{"id":"` + userId + `"},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"44e717f5d736b6d60427b7c1bb6d0462cd6dba65df8923a7fcd52c118199a18d"}}}`)

	req, err := http.NewRequest(http.MethodPost, "https://gql.twitch.tv/gql", bytes.NewBuffer(json))
	req.Header.Add("Client-Id", "kimne78kx3ncx6brgo4mv6wki5h1ko")
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	data, ok := utils.ParseJsonBody[LastUpdatedData](res.Body)

	if !ok || data.Data.User == nil {
		return LastUpdatedUser{}, false
	}

	return *data.Data.User, true
}

func GetUsersLastUpdate(userIds []string) ([]LastUpdatedUser, bool) {
	userIdsString, _ := json.Marshal(userIds)
	json := []byte(`{"operationName":"UsersLastUpdate","variables":{"ids":` + string(userIdsString) + `},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"d1c18e903d1d2b1989307bdf995014887b73e1415bebb19c63c0301927c118af"}}}`)

	req, err := http.NewRequest(http.MethodPost, "https://gql.twitch.tv/gql", bytes.NewBuffer(json))
	req.Header.Add("Client-Id", "kimne78kx3ncx6brgo4mv6wki5h1ko")
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	data, ok := utils.ParseJsonBody[LastUpdatedDataArray](res.Body)

	if !ok || len(data.Data.Users) == 0 {
		return []LastUpdatedUser{}, false
	}

	return data.Data.Users, true
}
