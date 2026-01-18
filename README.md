# moviefeed

RSS feed generator for tracking recent TV show episodes from TMDB. Fetches the latest episodes from your followed shows and serves them as an RSS feed.

## Quick Start

```bash
# Build
go build -o moviefeed .

# Create config
cp config.example.yaml config.yaml
# Edit config.yaml with your TMDB API key and show IDs

# Run
./moviefeed -config config.yaml
```

Server starts on port 8000 (configurable). Access the RSS feed at `http://localhost:8000/`, or `http://localhost:8000/?access_key=mySecretKey`

## Configuration

Create `config.yaml` (or `config.json`):

```yaml
api_key: "your_tmdb_api_key"  # required: get from themoviedb.org/settings/api
access_key: "sercret-key" # adds a layer of "auth", only users that know the key can access the api
port: "8000"  # optional: defaults to "8000"
shows:  # required
  - "tt22202452"  # IMDB ID (Pluribus)
  - "1396"       # TMDB ID (Breaking Bad)
```

## How It Works

- Fetches first season and latest season for each show
- Filters episodes aired in the last 30 days
- Returns episodes in reverse chronological order (newest first)
- Continues processing if individual shows fail to fetch
