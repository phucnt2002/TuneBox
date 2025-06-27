package controller

import (
	"TuneBox/domain"
	"context"
	"encoding/json"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"net/http"
)

type YoutubeController struct {
	service *youtube.Service
}

func NewYouTubeController(apiKey string) (*YoutubeController, error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &YoutubeController{service: service}, nil
}

func (c *YoutubeController) SearchSongs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}
	call := c.service.Search.List([]string{"snippet"}).
		Q(query).
		Type("video").
		MaxResults(5)
	resp, err := call.Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var result []domain.Song
	for _, item := range resp.Items {
		result = append(result, domain.Song{
			Title:   item.Snippet.Title,
			VideoID: item.Id.VideoId,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
