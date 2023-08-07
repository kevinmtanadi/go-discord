package service

import (
	"fmt"
	"go-discord/song"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

func PlaySong(s *discordgo.Session) {
	songChannelID := os.Getenv("SONG_CHANNEL_ID")

	for {
		songLists := song.GetSongListInstance()
		if len(songLists.Songs) == 0 {
			time.Sleep(1 * time.Second)
		}

		currentSong := songLists.Songs[0]

		var wg sync.WaitGroup
		wg.Add(1)

		go DownloadAudio(currentSong.URL, &wg)

		wg.Wait()

		embed := discordgo.MessageEmbed{
			Description: fmt.Sprintf("Currently playing **%s**", currentSong.Title),
		}
		s.ChannelMessageSendEmbed(songChannelID, &embed)
		time.Sleep(currentSong.Duration)

		songLists.Songs = songLists.Songs[:len(songLists.Songs)-1]
	}
}
