package repository

import (
	"TuneBox/domain"
	"sync"
)

type GroupRepository struct {
	group map[string]*domain.Group
	mutex sync.Mutex
}

func NewGroupRepository() *GroupRepository {
	return &GroupRepository{
		group: make(map[string]*domain.Group),
	}
}

func (r *GroupRepository) CreateGroup(name string) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	id := name
	if _, exits := r.group[id]; !exits {
		r.group[id] = &domain.Group{
			ID:       id,
			Name:     name,
			Playlist: make([]domain.Song, 0),
		}
	}
	return id
}
func (r *GroupRepository) GetGroup(id string) *domain.Group {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.group[id]
}

func (r *GroupRepository) AddSong(id string, song domain.Song) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if group, exits := r.group[id]; exits {
		group.Playlist = append(group.Playlist, song)
	}
}

func (r *GroupRepository) RemoveSong(groupID string, index int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if group, exites := r.group[groupID]; exites && index >= 0 && index < len(group.Playlist) {
		group.Playlist = append(group.Playlist[:index], group.Playlist[index+1:]...)
	}
}

func (r *GroupRepository) GetNextSong(groupID string) (domain.Song, []domain.Song) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if group, exites := r.group[groupID]; exites && len(group.Playlist) > 0 {
		song := group.Playlist[0]
		group.Playlist = group.Playlist[1:]
		return song, group.Playlist
	}
	return domain.Song{}, make([]domain.Song, 0)
}
