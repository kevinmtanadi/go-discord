package song

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type Song struct {
	Title              string
	Duration           time.Duration
	URL                string
	RequesterChannelID string
	VoiceChannelID     string
	Requester          *discordgo.User
}

type SongList struct {
	Songs map[string][]Song
}

var songListInstance *SongList

func init() {
	songListInstance = &SongList{
		Songs: make(map[string][]Song, 0),
	}
}

func GetSongListInstance() *SongList {
	return songListInstance
}

func (s *SongList) AddSong(song Song, guildID string) {
	s.Songs[guildID] = append(s.Songs[guildID], song)
}

func (s *SongList) Clear(guildID string) {
	s.Songs[guildID] = []Song{}
}
