package renamer

// LoginRequest is used to make login requests to tvdb api
type LoginRequest struct {
	Apikey   string
	Userkey  string
	Username string
}

// LoginResponse is used to map login responses from tvdb api
type LoginResponse struct {
	Token string
}

// SearchResponse is used to make search requests to tvdb api
type SearchResponse struct {
	Data []SeriesData
}

// SeriesData is used with search requests to tvdb api. Search requests return multiple
// SeriesData structs for a single search request.
type SeriesData struct {
	Aliases    []string
	FirstAired string
	ID         int
	SeriesName string
	Slug       string
	Status     string
}

// EpisodeResponse is used to map responses from tvdb api for episode details for a specific series
type EpisodeResponse struct {
	Links EpisodeListDetails
	Data  []EpisodeDetails
}

// EpisodeDetails contains information for a single episode of a series
type EpisodeDetails struct {
	AiredSeason        int
	AiredSeasonID      int
	AiredEpisodeNumber int
	EpisodeName        string
}

// EpisodeListDetails contains details for the number of pages of episode for a specific series
type EpisodeListDetails struct {
	First int
	Last  int
}

// ErrorResponse is used to map error responses from tvdb api
type ErrorResponse struct {
	Error string
}

// EpisodeNumber contains the episode and season number
type EpisodeNumber struct {
	Season  int
	Episode int
}
