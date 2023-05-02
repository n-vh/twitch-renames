package graphql

import (
	"time"
)

// GQL LAST UPDATED

type LastUpdatedData struct {
	Data       LastUpdatedUserData `json:"data"`
	Extensions Extensions          `json:"extensions"`
}

type LastUpdatedUserData struct {
	User *LastUpdatedUser `json:"user"`
}

type LastUpdatedUser struct {
	UserId      string    `json:"userId"`
	Login       string    `json:"login"`
	DisplayName string    `json:"displayName"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Extensions struct {
	DurationMilliseconds int64  `json:"durationMilliseconds"`
	OperationName        string `json:"operationName"`
	RequestID            string `json:"requestID"`
}

type LastUpdatedDataArray struct {
	Data       LastUpdatedUsersData `json:"data"`
	Extensions Extensions           `json:"extensions"`
}

type LastUpdatedUsersData struct {
	Users []LastUpdatedUser `json:"users"`
}
