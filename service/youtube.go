package service

import (
	"context"
	"errors"
	"fmt"
	"go-discord/song"
	"os"
	"regexp"
	"strconv"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// SearchYoutube searches for a song on YouTube based on the given query.
//
// Parameters:
//
//	query (string): The search query.
//
// Returns:
//
//	*song.Song: The song that was found.
//	error: An error if the search or retrieval of song details fails.
func SearchYoutube(query string) (*song.Song, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")

	// Create a new YouTube service
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	var videoID string
	if !isYouTubeLink(query) {
		searchResponse, err := service.Search.List([]string{"id,snippet"}).Q(query).MaxResults(1).Do()
		if err != nil {
			return nil, err
		}
		if len(searchResponse.Items) == 0 {
			return nil, errors.New("no videos found")
		}
		searchResult := searchResponse.Items[0]
		videoID = searchResult.Id.VideoId
	} else {
		videoID = getYouTubeVideoID(query)
	}

	video, err := service.Videos.List([]string{"snippet,contentDetails"}).Id(videoID).Do()
	if err != nil {
		return nil, err
	}
	if len(video.Items) == 0 {
		return nil, errors.New("video details not found")
	}

	videoStream := video.Items[0]
	duration, err := parseDuration(videoStream.ContentDetails.Duration)
	if err != nil {
		return nil, err
	}
	song := &song.Song{
		URL:      fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
		Title:    videoStream.Snippet.Title,
		Duration: duration,
	}

	return song, nil
}

func getYouTubeVideoID(link string) string {
	pattern := `^(https?://)?(www\.)?youtube\.com/watch\?v=([a-zA-Z0-9_-]{11})$`

	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(link)
	if len(matches) == 4 {
		return matches[3]
	}

	return ""
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

func isYouTubeLink(link string) bool {
	// Regular expression pattern to match a YouTube video URL
	pattern := `^(https?://)?(www\.)?youtube\.com/watch\?v=[a-zA-Z0-9_-]{11}$`

	match, _ := regexp.MatchString(pattern, link)
	return match
}
