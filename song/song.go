package song

import (
	"fmt"
	"time"
)

type Song struct {
	Title    string
	Duration time.Duration
	URL      string
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

func (s *SongList) PlaySong() {
	go func() {
		for {
			if len(s.Songs) > 0 {
				// PLAY THE SONG
				currentSong := s.Songs[0]
				fmt.Println(currentSong)

				time.Sleep(currentSong.Duration)
			}
		}
	}()
}
