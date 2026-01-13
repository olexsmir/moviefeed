package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type TMDBShow struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	FirstAirDate string `json:"first_air_date"`
}

type tmdbShowDetails struct {
	TMDBShow
	NumberOfSeasons int `json:"number_of_seasons"`
}

type TMDBEpisode struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Overview      string `json:"overview"`
	AirDate       string `json:"air_date"`
	EpisodeNumber int    `json:"episode_number"`
	SeasonNumber  int    `json:"season_number"`
	StillPath     string `json:"still_path"`
	ShowName      string
	ShowID        string
}

type tmdbFindResponse struct {
	TvResults []TMDBShow `json:"tv_results"`
}

type tmdbSeasonResponse struct {
	Episodes []TMDBEpisode `json:"episodes"`
}

func fetchNewEpisodes(config *Config) ([]TMDBEpisode, error) {
	var allEpisodes []TMDBEpisode
	for _, showID := range config.Shows {
		episodes, err := fetchEpisodesForShow(showID, config.APIKey)
		if err != nil {
			slog.Warn("failed to fetch episodes for show", "show", showID, "err", err)
			continue
		}
		allEpisodes = append(allEpisodes, episodes...)
	}
	return allEpisodes, nil
}

func getTMDBID(showID, apiKey string) (string, error) {
	imdbIDPrefix := "tt"
	if len(showID) > len(imdbIDPrefix) && showID[:len(imdbIDPrefix)] == imdbIDPrefix {
		url := fmt.Sprintf(
			"https://api.themoviedb.org/3/find/%s?api_key=%s&external_source=imdb_id",
			showID,
			apiKey,
		)
		result, err := makeRequest[tmdbFindResponse](url)
		if err != nil {
			return "", err
		}

		if len(result.TvResults) == 0 {
			return "", fmt.Errorf("no TMDB show found for IMDB ID %s", showID)
		}

		return fmt.Sprintf("%d", result.TvResults[0].ID), nil
	}
	return showID, nil
}

func fetchEpisodesForShow(showID, apiKey string) ([]TMDBEpisode, error) {
	tmdbID, err := getTMDBID(showID, apiKey)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s?api_key=%s", tmdbID, apiKey)
	show, err := makeRequest[tmdbShowDetails](url)
	if err != nil {
		return nil, err
	}

	seasonsToFetch := []int{1}
	if show.NumberOfSeasons > 1 {
		seasonsToFetch = append(seasonsToFetch, show.NumberOfSeasons)
	}

	var allEpisodes []TMDBEpisode
	for _, season := range seasonsToFetch {
		seasonURL := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s/season/%d?api_key=%s",
			tmdbID, season, apiKey)
		seasonData, err := makeRequest[tmdbSeasonResponse](seasonURL)
		if err != nil {
			slog.Warn("failed to fetch season", "season", season, "show", tmdbID, "err", err)
			continue
		}

		for _, ep := range seasonData.Episodes {
			ep.ShowName = show.Name
			ep.ShowID = tmdbID
			allEpisodes = append(allEpisodes, ep)
		}
	}

	return filterRecentEpisodes(allEpisodes), nil
}

func filterRecentEpisodes(episodes []TMDBEpisode) []TMDBEpisode {
	var recent []TMDBEpisode
	now := time.Now()
	cutoff := now.AddDate(0, 0, -30)

	for _, ep := range episodes {
		if ep.AirDate == "" {
			continue
		}

		airDate, err := time.Parse(dateFormat, ep.AirDate)
		if err != nil {
			continue
		}

		if airDate.Before(now) && airDate.After(cutoff) {
			recent = append(recent, ep)
		}
	}
	return recent
}

func makeRequest[T any](url string) (*T, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "err", err)
		}
	}()

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
