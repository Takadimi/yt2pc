package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Takadimi/yt2pc/core/youtube"
)

func main() {
	// gets video data from innertube API for test video unless otherwise specified:
	// https://www.youtube.com/watch?v=BaW_jenozKc&pp=ygUVeW91dHViZS1kbCB0ZXN0IHZpZGVv
	videoID := "BaW_jenozKc"
	if len(os.Args) > 1 {
		videoID = os.Args[1]
	}

	videoData, err := youtube.VideoDataByInnertube(http.DefaultClient, videoID)
	if err != nil {
		log.Fatal(err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(videoData, &jsonData); err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))
}
