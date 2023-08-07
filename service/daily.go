package service

import (
	"encoding/json"
	"go-discord/helper"
	"go-discord/logger"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

type Holiday struct {
	HolidayDate       string `json:"holiday_date"`
	HolidayName       string `json:"holiday_name"`
	IsNationalHoliday bool   `json:"is_national_holiday"`
}

func (h *Holiday) FormatString() string {
	result := ""
	result += "**" + h.HolidayName + "**\n"
	result += h.HolidayDate + "\n"

	return result
}

func GetHoliday(s *discordgo.Session) {
	holidayChannelID := os.Getenv("HOLIDAYT_CHANNEL_ID")

	resp, err := http.Get(os.Getenv("HOLIDAY_API") + "api")
	if err != nil {
		logger.Log("Fail calling holiday API: " + err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Log("Fail reading data: " + err.Error())
		return
	}

	holidays := []Holiday{}
	err = json.Unmarshal(body, &holidays)
	if err != nil {
		logger.Log("Fail binding data: " + err.Error())
		return
	}

	closeHolidays := ""
	for _, holiday := range holidays {
		if !holiday.IsNationalHoliday {
			continue
		}

		date, err := helper.ParseDate(holiday.HolidayDate)
		if err != nil {
			logger.Log("Fail parsing date: " + err.Error())
		}

		currentTime := time.Now()
		threshold := currentTime.AddDate(0, 0, 30)

		if !date.After(currentTime) || !date.Before(threshold) {
			continue
		}
		holiday.HolidayDate = date.Format("02-01-2006")

		closeHolidays += holiday.FormatString()
	}

	if len(closeHolidays) == 0 {
		embed := discordgo.MessageEmbed{
			Title:       "Hari Libur 30 hari kedepan",
			Description: "Saat ini tidak ada libur",
		}
		s.ChannelMessageSendEmbed(holidayChannelID, &embed)

		return
	}

	embed := discordgo.MessageEmbed{
		Title:       "Hari Libur 30 hari kedepan",
		Description: closeHolidays,
	}
	s.ChannelMessageSendEmbed(holidayChannelID, &embed)
}

func DailyCall(s *discordgo.Session) {
	c := cron.New()

	c.AddFunc("0 8 0 * * *", func() {
		GetHoliday(s)
	})

	c.Start()
}
