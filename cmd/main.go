package main

import (
	"TuneBox/api/route"
	"TuneBox/bootstrap"
	"log"
	"net/http"
)

func main() {
	// Fix: assign both env and err
	env := bootstrap.NewEnv()
	router := &route.Route{
		Env: env,
	}

	// Fix: only one value being assigned here
	mux := route.SetupRoutes(router)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
