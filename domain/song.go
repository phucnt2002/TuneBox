package domain

type Song struct {
	Title   string `json:"title"`
	VideoID string `json:"videoID"`
}

type Group struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Playlist []Song `json:"playlist"`
}
