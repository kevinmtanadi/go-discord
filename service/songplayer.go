package service

import (
	"fmt"
	"go-discord/helper"
	"go-discord/logger"
	"go-discord/song"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func PlaySong(s *discordgo.Session) {
	guildId := os.Getenv("GUILD_ID")
	songChannelID := os.Getenv("VOICE_CHANNEL_ID")

	var wg sync.WaitGroup
	for {
		helper.DeleteFileExists("audio.webm")
		songLists := song.GetSongListInstance()
		if len(songLists.Songs) == 0 {
			time.Sleep(1 * time.Second)
			voiceConn := s.VoiceConnections
			if _, ok := voiceConn[guildId]; ok {
				voiceConn[guildId].Disconnect()
			}
			continue
		}

		currentSong := songLists.Songs[0]

		wg.Add(1)

		go DownloadAudio(currentSong.URL, &wg)

		wg.Wait()
		time.Sleep(3 * time.Second)

		voice, err := s.ChannelVoiceJoin(guildId, songChannelID, false, true)
		if err != nil {
			logger.Log("Error joining voice channel: " + err.Error())
		}

		embed := discordgo.MessageEmbed{
			Description: fmt.Sprintf("Currently playing **%s**", currentSong.Title),
		}
		s.ChannelMessageSendEmbed(currentSong.ChannelID, &embed)

		s.UpdateGameStatus(0, currentSong.Title)
		dgvoice.PlayAudioFile(voice, "audio.webm", make(chan bool))
		fmt.Println(currentSong.Title)

		songLists.Songs = songLists.Songs[1:]
	}
}
