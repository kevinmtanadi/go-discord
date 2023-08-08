package service

import (
	"fmt"
	"go-discord/logger"
	"sync"

	"github.com/iawia002/lux/downloader"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/extractors/youtube"
)

// DownloadAudio downloads audio from the given URL.
//
// Parameters:
// - url: a string representing the URL from which to download the audio.
// - wg: a pointer to a sync.WaitGroup that is used to synchronize multiple goroutines.
//
// Return type: None.

var streamOptions = []string{"249", "250", "251"}

func DownloadAudio(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	filename := "audio"

	e := youtube.New()
	data, err := e.Extract(url, extractors.Options{})
	if err != nil {
		logger.Log("Fail extracting: " + err.Error())
		return
	}

	result := data[0]

	i := 0
	maxTry := 50

	for {
		if i >= maxTry {
			logger.Log(fmt.Sprintf("Failed to download %s after %d tries", result.Title, maxTry))
		}
		// Keep trying to download
		curOpt := streamOptions[i%len(streamOptions)]
		options := downloader.Options{
			Silent:         true,
			InfoOnly:       false,
			Stream:         curOpt,
			OutputPath:     "./",
			OutputName:     filename,
			FileNameLength: 64,
			Caption:        false,
			MultiThread:    true,
			ThreadNumber:   8,
			RetryTimes:     10,
		}
		defaultDownloader := downloader.New(options)

		err = defaultDownloader.Download(result)
		if err == nil {
			return
		}
		i++
	}
}
