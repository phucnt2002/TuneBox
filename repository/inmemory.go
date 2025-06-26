package repository

import (
	"TuneBox/domain"
	"sync"
)

type InMemoryRepository struct {
	playList []domain.Song
	mutex    sync.Mutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		playList: make([]domain.Song, 0)}
}

func (r *InMemoryRepository) AddSong(song domain.Song) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.playList = append(r.playList, song)
}

func (r *InMemoryRepository) RemoveSong(index int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if index >= 0 && index < len(r.playList) {
		r.playList = append(r.playList[:index], r.playList[index+1:]...)
	}
}

func (r *InMemoryRepository) GetNextSong() (domain.Song, []domain.Song) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if len(r.playList) == 0 {
		return domain.Song{}, r.playList
	}
	song := r.playList[0]
	r.playList = r.playList[1:]
	return song, r.playList
}

func (r *InMemoryRepository) GetPlayList() []domain.Song {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.playList
}
