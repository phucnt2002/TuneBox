package bootstrap

import "os"

func GetAPIKey() string {
	return os.Getenv("YOUTUBE_API_KEY")
}
