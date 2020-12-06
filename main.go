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
)

var (
	configs config.Configurations
	db database.Database = nil
)

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
	videoInference, err := db.GetVideoInference(r.Context(), ps.ByName("name"))
	if err != nil {
		fmt.Fprintf(w, "Unable to get ad information")
		return
	}
	log.Println("Fetched video inference")

	keys := make([]string, 0, len(videoInference.TopFiveObjectsToInterval))
	for k := range videoInference.TopFiveObjectsToInterval {
		keys = append(keys, k)
	}

	ads, err := db.FindAdsWithObjects(r.Context(), keys)
	if err != nil {
		fmt.Fprintf(w, "Unable to get ad information")
		return
	}
	log.Println("Fetched relevant ads")

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(ads)
	log.Println("Sent relevant ad information")
}

func init() {
	var err error
	configs, err = config.GetConfigs()
	if err != nil {
		log.Println("Unable to get config")
	}
}

func main() {
	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile )

	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// DB and store
	var err error
	db, err = database.GetDatabaseClient(ctx, configs.Database)
	if err != nil {
		log.Fatalln(err)
	}

	router := httprouter.New()
	router.GET("/videos", videoGetAll)
	router.GET("/videos/:name", getVideoAndAdInfo)

	log.Println("Serving on", configs.Server.Port)
	http.ListenAndServe(":" + configs.Server.Port, router)
}