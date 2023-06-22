package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Takadimi/yt2pc/core/youtube"
)

func main() {
	// gets video data from innertube API for test video:
	// https://www.youtube.com/watch?v=BaW_jenozKc&pp=ygUVeW91dHViZS1kbCB0ZXN0IHZpZGVv
	videoData, err := youtube.VideoDataByInnertube(http.DefaultClient, "BaW_jenozKc")
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
