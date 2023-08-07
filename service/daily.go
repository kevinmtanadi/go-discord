package service

import (
	"go-discord/helper"
	"go-discord/logger"
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

// GetHoliday retrieves upcoming holidays and sends a message to a Discord channel.
// It takes a *discordgo.Session pointer as a parameter and does not return anything.
func GetHoliday(s *discordgo.Session) {
	holidayChannelID := os.Getenv("HOLIDAY_CHANNEL_ID")

	holidays := []Holiday{}
	err := helper.GetResponse(os.Getenv("HOLIDAY_API")+"api", &holidays)
	if err != nil {
		logger.Log("Fail getting data: " + err.Error())
	}

	closeHolidays := ""
	currentTime := time.Now()
	threshold := currentTime.AddDate(0, 0, 30)

	for _, holiday := range holidays {
		if !holiday.IsNationalHoliday {
			continue
		}

		date, err := helper.ParseDate(holiday.HolidayDate)
		if err != nil {
			logger.Log("Fail parsing date: " + err.Error())
		}

		if !date.After(currentTime) || !date.Before(threshold) {
			continue
		}
		holiday.HolidayDate = date.Format("02-01-2006")

		closeHolidays += holiday.FormatString()
	}

	var embed discordgo.MessageEmbed
	if len(closeHolidays) == 0 {
		embed = discordgo.MessageEmbed{
			Title:       "Hari Libur 30 hari kedepan",
			Description: "Saat ini tidak ada libur",
		}
	} else {
		embed = discordgo.MessageEmbed{
			Title:       "Hari Libur 30 hari kedepan",
			Description: closeHolidays,
		}
	}

	previousEmbeds, err := s.ChannelMessages(holidayChannelID, 1, "", "", "")
	if err != nil {
		logger.Log("Fail getting previous messages: " + err.Error())
		return
	}

	if len(previousEmbeds) > 0 {
		err = s.ChannelMessageDelete(holidayChannelID, previousEmbeds[0].ID)
		if err != nil {
			logger.Log("Fail deleting previous message: " + err.Error())
			return
		}
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
