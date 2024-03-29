package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	database "github.com/vp-cap/data-lib/database"
	config "github.com/vp-cap/video-service/config"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var (
	configs config.Configurations
	db database.Database = nil
)

type videoAndAdInfo struct {
	Video database.Video           `json:"video"`
	Intervals  map[string]database.Interval  `json:"intervals"`
	Ads   []database.Advertisement `json:"ads"`
}

// handle video fetch all
func videoGetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	videos, err := db.GetAllVideos(r.Context())
	if err != nil {
		fmt.Fprintf(w, "Unable to fetch video information")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videos)
	log.Println("Sent all video information")
}

func getVideoAndAdInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	video, err := db.GetVideo(r.Context(), ps.ByName("name"))
	if err != nil {
		fmt.Fprintf(w, "Unable to get video information")
		return
	}
	log.Println("Fetched video")

	videoInference, err := db.GetVideoInference(r.Context(), video.StorageLink)
	if err != nil {
		log.Println("No Video inference found", err)
	} else {
		log.Println("Fetched video inference")
	}

	keys := make([]string, 0, len(videoInference.TopFiveObjectsToInterval))
	for k := range videoInference.TopFiveObjectsToInterval {
		keys = append(keys, k)
	}

	ads, err := db.FindAdsWithObjects(r.Context(), keys)
	if err != nil {
		log.Println("No Ads found", err)
	} else {
		log.Println("Fetched relevant ads")
	}
	
	videoAndAdInfo, err := json.Marshal(&videoAndAdInfo{Video: video, Intervals: videoInference.TopFiveObjectsToInterval, Ads: ads})
	if err != nil {
		fmt.Fprintf(w, "Unable to marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(videoAndAdInfo))

	log.Println("Sent video and relevant ad information")
}

func init() {
	var err error
	configs, err = config.GetConfigs()
	if err != nil {
		log.Fatal("Unable to get config")
	}
}

func main() {
	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile )

	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	var err error
	db, err = database.GetDatabaseClient(ctx, configs.Database)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to DB")

	router := httprouter.New()
	router.GET("/videos", videoGetAll)
	router.GET("/videos/:name", getVideoAndAdInfo)

	handler := cors.Default().Handler(router)

	log.Println("Serving on", configs.Server.Port)
	http.ListenAndServe(":" + configs.Server.Port, handler)
}