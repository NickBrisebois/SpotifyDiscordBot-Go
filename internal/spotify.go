package internal

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	auth  spotify.Authenticator
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func InitSpotify(config *Config) (err error) {
	//	startAuthServer(config)

	spotifyConfig := &clientcredentials.Config{
		ClientID:     config.SpotifyClientID,
		ClientSecret: config.SpotifySecretKey,
		TokenURL:     spotify.TokenURL,
	}

	token, err := spotifyConfig.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	auth := spotify.NewAuthenticator(config.SpotifyRedirectURL, spotify.ScopeUserReadPrivate)
	auth.SetAuthInfo(config.SpotifyClientID, config.SpotifySecretKey)

	client := auth.NewClient(token)

	const PlaylistID spotify.ID = "37i9dQZF1EthtctLd3ak1i"
	results, err := client.GetPlaylist(PlaylistID)
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	fmt.Println(results.Name)

	// Listen for signals so that the program is capable of quitting
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return nil
}

/*
func startAuthServer(config *Config) {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})

	go http.ListenAndServe(":8080", nil)

	fmt.Println(config)
	auth = spotify.NewAuthenticator(config.SpotifyRedirectURL, spotify.ScopeUserReadPrivate)
	url := auth.AuthURL(state)
	fmt.Println("Please login to Spotiify by visiting the following page in your browser: ", url)

	client := <-ch

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Logged in as: ", user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		log.Println(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s", st, state)
	}

	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
*/
