package internal

// Config holds spotify and discord configuration data
type Config struct {
	DiscordToken        string
	SpotifyRedirectURL  string
	SpotifyClientID     string
	SpotifyClientSecret string
	LimitToOneChannel   bool
	ChannelToUse        string
	SpotifyPlaylist     string
}
