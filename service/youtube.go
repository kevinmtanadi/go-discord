package service

import (
	"bytes"
	"fmt"
	"go-discord/helper"
	"go-discord/song"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// TODO USE A MORE EFFICIENT YOUTUBE SEARCH
func Searchyoutube(query string, guildID string) (*song.Song, error) {
	fileName := guildID + ".webm"
	var title string
	var durationStr string

	if !isYouTubeLink(query) {
		output, err := exec.Command("./youtube-dl", "--get-title", "--get-duration", "ytsearch1:"+query).Output()
		if err != nil {
			return nil, err
		}

		bufs := new(bytes.Buffer)
		wr := transform.NewWriter(bufs, japanese.ShiftJIS.NewDecoder())
		wr.Write(output)
		wr.Close()

		result := strings.Split(string(output), "\n")
		title = result[0]
		durationStr = result[1]
	} else {
		cmd := exec.Command("./youtube-dl", "--get-title", "--get-duration", query)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return nil, err
		}

		result := strings.Split(string(output), "\n")
		title = result[0]
		durationStr = result[1]
	}

	duration := helper.ParseStringToTimeDuration(durationStr)

	songSearched := song.Song{
		Title:    title,
		Duration: duration,
		Filename: fileName,
	}

	return &songSearched, nil
}

func DownloadYoutube(song *song.Song, guildID string) error {
	if !isYouTubeLink(song.SearchQuery) {
		// exec youtube-dl command
		cmd := exec.Command("./youtube-dl", "-f", "250", "-o", song.Filename, "-q", "ytsearch1:"+song.SearchQuery)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return err
		}

	} else {
		cmd := exec.Command("./youtube-dl", "-f", "250", "-o", song.Filename, "-q", song.SearchQuery)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return err
		}
	}

	return nil

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

func ScrapeYoutube(query string) {
	encodedQuery := url.QueryEscape(query)
	searchURL := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", encodedQuery)

	response, err := http.Get(searchURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Parse the HTML response using goquery
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	html, err := ioutil.ReadAll(response.Body)

	htmlText := string(html)
	regexPattern := `/(?<=shortDescription":").*(?=","isCrawlable)/`

	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(htmlText)
	if len(matches) > 0 {
		fmt.Println(matches)
	} else {
		fmt.Println("No matches found")
	}

	fmt.Println("Scrape result: ")
	var firstVideoLink string
	doc.Find(`(?<=shortDescription":").*(?=","isCrawlable)`).Each(func(index int, item *goquery.Selection) {
		if index == 0 {
			// Extract the video link
			videoLink, _ := item.Attr("href")
			firstVideoLink = "https://www.youtube.com" + videoLink
			return
		}
	})
	fmt.Println(firstVideoLink)
	fmt.Println("EOL")
}

func TestScrape() {
	html := `
        <html>
            <head>
                <title>Sample HTML Document</title>
            </head>
            <body>
                <h1>Hello, goquery!</h1>
                <p>Welcome to the world of web scraping.</p>
            </body>
        </html>
    `

	reader := strings.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal("Error creating document:", err)
	}

	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
}
