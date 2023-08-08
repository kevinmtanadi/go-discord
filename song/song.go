package song

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Song struct {
	Title     string
	Duration  time.Duration
	URL       string
	ChannelID string
	Requester *discordgo.User
}

type SongList struct {
	Songs []Song
}

var songListInstance *SongList

func init() {
	songListInstance = &SongList{
		Songs: make([]Song, 0),
	}
}

func GetSongListInstance() *SongList {
	return songListInstance
}

func (s *SongList) AddSong(song Song) {
	s.Songs = append(s.Songs, song)
}

func (s *SongList) Clear() {
	s.Songs = []Song{}
}
