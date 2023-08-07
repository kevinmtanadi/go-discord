package service

import (
	"errors"
	"fmt"
	"go-discord/song"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// SearchYoutube searches for a song on YouTube based on the given query.
//
// Parameters:
//   query (string): The search query.
//
// Returns:
//   *song.Song: The song that was found.
//   error: An error if the search or retrieval of song details fails.
func SearchYoutube(query string) (*song.Song, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")

	httpClient := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}

	service, err := youtube.New(httpClient)
	if err != nil {
		return nil, err
	}

	searchResponse, err := service.Search.List([]string{"id,snippet"}).Q(query).MaxResults(1).Do()
	if err != nil {
		return nil, err
	}
	if len(searchResponse.Items) == 0 {
		return nil, errors.New("no videos found")
	}

	video := searchResponse.Items[0]
	videoID := video.Id.VideoId
	videoResponse, err := service.Videos.List([]string{"snippet,contentDetails"}).Id(videoID).Do()
	if err != nil {
		return nil, err
	}
	if len(videoResponse.Items) == 0 {
		return nil, errors.New("video details not found")
	}

	duration, err := parseDuration(videoResponse.Items[0].ContentDetails.Duration)
	if err != nil {
		return nil, err
	}
	song := &song.Song{
		URL:      fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		Title:    video.Snippet.Title,
		Duration: duration,
	}

	return song, nil
}

func parseDuration(durationString string) (time.Duration, error) {
	re := regexp.MustCompile(`PT(\d+)M(\d+)S`)

	// Find matches using the regular expression.
	matches := re.FindStringSubmatch(durationString)

	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid duration format")
	}

	// Extract minutes and seconds from the matched groups.
	minutes, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, err
	}

	// Calculate the total duration in seconds.
	totalSeconds := minutes*60 + seconds

	// Convert to a time.Duration value.
	duration := time.Duration(totalSeconds) * time.Second

	return duration, nil
}
