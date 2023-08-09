package service

import (
	"fmt"
	"go-discord/helper"
	"go-discord/logger"
	"go-discord/song"
	"sync"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func PlaySong(s *discordgo.Session) {
	var wg sync.WaitGroup
	for {
		songLists := song.GetSongListInstance()
		for guildID, songs := range songLists.Songs {
			helper.DeleteFileExists(guildID + ".webm")
			if len(songs) == 0 {
				time.Sleep(1 * time.Second)
				voiceConn := s.VoiceConnections
				if _, ok := voiceConn[guildID]; ok {
					voiceConn[guildID].Disconnect()
				}
				continue
			}
			currentSong := songs[0]

			wg.Add(1)

			go DownloadAudio(currentSong.URL, &wg, guildID)

			wg.Wait()
			time.Sleep(3 * time.Second)

			voice, err := s.ChannelVoiceJoin(guildID, currentSong.VoiceChannelID, false, true)
			if err != nil {
				logger.Log("Error joining voice channel: " + err.Error())
			}

			embed := discordgo.MessageEmbed{
				Description: fmt.Sprintf("Currently playing **%s**", currentSong.Title),
			}
			s.ChannelMessageSendEmbed(currentSong.RequesterChannelID, &embed)

			s.UpdateGameStatus(0, currentSong.Title)
			dgvoice.PlayAudioFile(voice, guildID+".webm", make(chan bool))
			helper.DeleteFileExists(guildID + ".webm")

			songLists.Songs[guildID] = songs[1:]
		}

	}
}
