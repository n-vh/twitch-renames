package helix

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/n-vh/twitch-renames/internal/utils"
	"github.com/samber/lo"
)

func GetStreams(cursor string) (Streams, bool) {
	params := url.Values{}
	params.Add("language", "en")
	params.Add("language", "fr")
	params.Add("language", "es")
	params.Add("language", "pt")
	params.Add("language", "de")
	params.Add("first", "100")

	if len(cursor) != 0 {
		params.Add("after", cursor)
	}

	url := fmt.Sprintf("%s/streams?%s", ENDPOINT, params.Encode())
	res, err := http.Get(url)

	if err != nil {
		log.Panic(err)
	}

	streams := Streams{}

	data, ok := utils.ParseJsonBody[HelixStreamData](res.Body)

	if !ok {
		return streams, false
	}

	streams.Cursor = data.Pagination.Cursor

	for _, stream := range data.Data {
		streams.Streams = append(streams.Streams, Stream(stream))
	}

	return streams, true
}

func GetStreamsPaginated(total int) []Stream {
	log.Println("Fetching streams")

	var streams []Stream

	for cursor := ""; len(streams) <= total; {
		data, ok := GetStreams(cursor)

		if !ok || len(data.Streams) == 0 {
			break
		}

		cursor = data.Cursor

		for _, stream := range data.Streams {
			contains := lo.ContainsBy(streams, func(s Stream) bool {
				return s.UserId == stream.UserId
			})

			if contains {
				continue
			}

			streams = append(streams, stream)
		}
	}

	return lo.Slice(streams, 0, total)
}
