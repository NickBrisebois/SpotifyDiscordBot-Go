package spotifyauth

import (
	"fmt"
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/config"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

var (
	authenticator spotify.Authenticator
	ch            = make(chan *spotify.Client)
	state         = "asdfasdfasdfasdfasdfasdfa"
)

// GetSpotifyClient handles authentication and returns a spotify client
func GetSpotifyClient(spottyConf *config.Config) *spotify.Client {
	spottyScopes := []string{
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistReadPrivate,
		spotify.ScopeUserLibraryRead,
	}

	authenticator = spotify.NewAuthenticator(spottyConf.SpotifyRedirectURL, spottyScopes...)
	authenticator.SetAuthInfo(spottyConf.SpotifyClientID, spottyConf.SpotifyClientSecret)

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := authenticator.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	return client
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := authenticator.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := authenticator.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
