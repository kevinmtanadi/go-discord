package service

import (
	"fmt"
	"go-discord/helper"
	"go-discord/logger"
	"sync"

	"github.com/iawia002/lux/downloader"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/extractors/youtube"
)

func DownloadAudio(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	filename := "audio"
	helper.DeleteFileExists(filename + ".webm")

	e := youtube.New()
	data, err := e.Extract(url, extractors.Options{})

	if err != nil {
		logger.Log("Fail extracting: " + err.Error())
	}

	result := data[0]
	fmt.Println(result.Streams)

	defaultDownloader := downloader.New(downloader.Options{
		Silent:         false,
		InfoOnly:       false,
		Stream:         "250",
		OutputPath:     "./",
		OutputName:     filename,
		FileNameLength: 10,
		Caption:        false,
		MultiThread:    true,
		ThreadNumber:   8,
		RetryTimes:     1,
	})

	err = defaultDownloader.Download(result)
	if err != nil {
		logger.Log("Fail extracting: " + err.Error())
	}

}
