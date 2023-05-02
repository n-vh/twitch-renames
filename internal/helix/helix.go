package helix

const (
	ENDPOINT = "http://localhost:8535/helix"
)

// HELIX STREAM

type HelixStreamData struct {
	Data       []HelixStream `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

type HelixStream struct {
	UserId      string `json:"user_id"`
	Login       string `json:"user_login"`
	DisplayName string `json:"user_name"`
	Total       int    `json:"viewer_count"`
}

type Pagination struct {
	Cursor string `json:"cursor,omitempty"`
}

type Streams struct {
	Streams []Stream
	Cursor  string
}

type Stream struct {
	UserId      string
	Login       string
	DisplayName string
	Total       int
}

// HELIX USER

type HelixUserData struct {
	Data []HelixUser `json:"data"`
}

type HelixUser struct {
	UserId      string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"display_name"`
}

type User struct {
	UserId      string
	Login       string
	DisplayName string
}
