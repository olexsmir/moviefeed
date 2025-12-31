package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/feeds"
)

const dateFormat = "2006-01-02"

func main() {
	configFile := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	config, err := loadConfig(*configFile)
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		episodes, err := fetchNewEpisodes(config)
		if err != nil {
			slog.Error("failed to fetch episodes", "err", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		rssFeed := generateRSS(episodes)
		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(rssFeed))
	})

	addr := ":" + config.Port
	slog.Info("server starting", "port", config.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}

func generateRSS(episodes []TMDBEpisode) string {
	feed := &feeds.Feed{
		Title:       "moviefeed",
		Link:        &feeds.Link{Href: "http://localhost"},
		Description: "Latest episodes from followed shows",
		Created:     time.Now(),
	}

	for i := len(episodes) - 1; i >= 0; i-- {
		ep := episodes[i]
		airDate, _ := time.Parse("2006-01-02", ep.AirDate)
		feed.Items = append(feed.Items, &feeds.Item{
			Id: fmt.Sprintf("%s-%d-%d", ep.ShowID, ep.SeasonNumber, ep.EpisodeNumber),
			Title: fmt.Sprintf(
				"%s S%dE%d: %s",
				ep.ShowName,
				ep.SeasonNumber,
				ep.EpisodeNumber,
				ep.Name,
			),
			Link: &feeds.Link{
				Href: fmt.Sprintf("https://www.themoviedb.org/tv/episode/%d", ep.ID),
			},
			Description: ep.Overview,
			Created:     airDate,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		slog.Error("error generating RSS", "err", err)
		return ""
	}

	return rss
}
