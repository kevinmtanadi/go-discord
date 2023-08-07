package logger

import (
	"log"
	"os"
	"strings"
)

// Log : Log the message to stdout and log files
func Log(message string) {
	absolutePath := os.Getenv("LOG_PATH")
	path := strings.Split(absolutePath, "/")
	folderPath, filePath := path[0], path[1]
	CreateFile(folderPath, filePath)

	file, err := os.OpenFile(absolutePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println("=== " + message)
}

// CreateFile :  Create a folder and file if not exist
func CreateFile(folderPath, filePath string) {
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(folderPath + "/" + filePath); os.IsNotExist(err) {
		file, err := os.Create(folderPath + "/" + filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
}
