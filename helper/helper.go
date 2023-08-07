package helper

import (
	"encoding/json"
	"fmt"
	"go-discord/logger"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// func ConvertInterfaceToString(is ...interface{}) string {
// 	flattenedIntf := flattenInterfaceToString(is)
// 	res := []string{}
// 	for _, i := range is {
// 		fmt.Println(i)
// 		res = append(res, i.(string))
// 	}
// 	return strings.Join(res, " ")
// }

func GetArgs(input interface{}) string {
	args := ConvertInterfaceToString(input)
	return args[1 : len(args)-3]
}

func ConvertInterfaceToString(input interface{}) string {
	switch v := input.(type) {
	case string:
		return v
	case []interface{}:
		var result string
		for _, item := range v {
			result += ConvertInterfaceToString(item) + " "
		}
		return result
	default:
		return fmt.Sprint(v)
	}
}

func ParseDate(dateString string) (time.Time, error) {
	layout := "2006-01-02"

	// Split the date string by "-" separator
	parts := strings.Split(dateString, "-")

	// Extract year, month, and day from the parts
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, err
	}

	// Generate the formatted date string
	formattedDate := fmt.Sprintf("%04d-%02d-%02d", year, month, day)

	// Parse the formatted date string into a time.Time value
	date, err := time.Parse(layout, formattedDate)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func GetResponse(url string, result any) (err error) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}

	return nil
}

func DeleteFileExists(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		// File exists, delete it
		err := os.Remove(fileName)
		if err != nil {
			// Handle the error if the file deletion fails
			logger.Log("Failed deleting file: " + err.Error())
		}
	}
}
