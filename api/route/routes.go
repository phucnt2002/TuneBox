package route

import (
	"TuneBox/api/controller"
	"TuneBox/bootstrap"
	"TuneBox/repository"
	"log"
	"net/http"
)

type Route struct {
	Env *bootstrap.Env
}

func SetupRoutes(r *Route) *http.ServeMux {
	mux := http.NewServeMux()
	repo := repository.NewInMemoryRepository()
	wsController := controller.NewWebSocketController(repo)
	youtubeController, err := controller.NewYouTubeController(r.Env.YoutubeAPIKey)
	if err != nil {
		log.Fatalf("Error initializing YouTube controller: %v", err)
	}

	mux.HandleFunc("/ws", wsController.HandleConnections)
	mux.HandleFunc("/search", youtubeController.SearchSongs)
	mux.Handle("/", http.FileServer(http.Dir("./client/dist")))
	return mux
}
