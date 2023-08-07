package song

type Song struct {
	Title    string
	Duration string
	URL      string
}

type SongList struct {
	Songs []Song
}

func NewSongList() *SongList {
	return &SongList{
		Songs: make([]Song, 0),
	}
}

var songList SongList

func init() {
	songList = *NewSongList()
}

func AddSong(song Song) {
	songList.Songs = append(songList.Songs, song)
}

func (m *Song) Play() {

}
