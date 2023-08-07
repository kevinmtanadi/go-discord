package service

import (
	"fmt"
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

	for {
		songLists := song.GetSongListInstance()
		if len(songLists.Songs) == 0 {
			time.Sleep(1 * time.Second)
			voiceConn := s.VoiceConnections
			if _, ok := voiceConn[guildId]; ok {
				voiceConn[guildId].Close()
			}
			continue
		}

		currentSong := songLists.Songs[0]

		var wg sync.WaitGroup
		wg.Add(1)

		go DownloadAudio(currentSong.URL, &wg)

		wg.Wait()

		voice, err := s.ChannelVoiceJoin(guildId, songChannelID, false, true)
		if err != nil {
			logger.Log("Error joining voice channel: " + err.Error())
		}

		embed := discordgo.MessageEmbed{
			Description: fmt.Sprintf("Currently playing **%s**", currentSong.Title),
		}
		s.ChannelMessageSendEmbed(songChannelID, &embed)

		audioPath := "./audio.webm"
		s.UpdateGameStatus(0, currentSong.Title)
		dgvoice.PlayAudioFile(voice, audioPath, make(chan bool))

		songLists.Songs = songLists.Songs[:len(songLists.Songs)-1]
	}
}
