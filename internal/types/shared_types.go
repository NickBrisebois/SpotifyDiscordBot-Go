package types

// SpottyMsgKind is an enum declaring what kind of msg is being passed between the spotify and bot thread
type SpottyMsgKind int

const (
	// Track represents a message to the spotify thread containing info on a new track
	Track SpottyMsgKind = iota
	// Playlist represents a message to the spotify thread containing info on a new playlist
	Playlist
	// Album represents a message to the spotify thread containing info on a new playlist
	Album
	// AddedBy represents a request for the name of the user who added the given track ID
	AddedBy
	// Error represents a response from the spotify thread indicating something went wrong
	Error
	// Response is a non-error response from the spotify thread
	Response
)

// SpottyMessage is passed between main thread and spotify thread as a form of communication
type SpottyMessage struct {
	Kind SpottyMsgKind // Kind identifies what kind of information is being passed
	ID   string        // ID of track, album or playlist to send to spotify thread
	Msg  string        // Message response from thread to other thread
	User string        // User who added song/album/playlist
	Err  error         // Self explanatory, something happened so this won't be nil and contain an error
}
