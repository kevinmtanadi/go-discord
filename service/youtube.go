package service

import (
	"errors"
	"fmt"
	"go-discord/song"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func SearchYoutube(query string) (*song.Song, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")

	httpClient := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}

	song := song.Song{}

	service, err := youtube.New(httpClient)
	if err != nil {
		return &song, err
	}

	searchResponse, err := service.Search.List([]string{"id,snippet"}).Q(query).MaxResults(1).Do()
	if err != nil {
		return &song, err
	}
	if len(searchResponse.Items) == 0 {
		return &song, errors.New("no videos found")
	}

	video := searchResponse.Items[0]
	videoID := video.Id.VideoId
	videoResponse, err := service.Videos.List([]string{"snippet,contentDetails"}).Id(videoID).Do()
	if err != nil {
		return &song, err
	}
	if len(videoResponse.Items) == 0 {
		return &song, errors.New("video details not found")
	}

	videoDuration := videoResponse.Items[0].ContentDetails.Duration

	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	song.URL = videoURL
	song.Title = video.Snippet.Title
	song.Duration = videoDuration
	return &song, nil
}
