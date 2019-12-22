package bot

import (
	"context"
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/types"
	"github.com/zmb3/spotify"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type database struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// SpotifyHandler handles spotify interactiosn
type SpotifyHandler struct {
	db           *database
	config       *types.Config
	playlistID   spotify.ID
	spottyClient *spotify.Client
	spottyChan   chan types.SpottyMessage
}

// SongRecord is what is saved to mongodb to keep track of what songs have been added
type SongRecord struct {
	ID      string
	AddedBy string
}

// NewSpotifyHandler creates a new spotify handler
func NewSpotifyHandler(config *types.Config, spottyChan chan types.SpottyMessage, client *spotify.Client) *SpotifyHandler {
	spotifyHandler := &SpotifyHandler{
		db:           newDatabase(config),
		config:       config,
		spottyClient: client,
		playlistID:   spotify.ID(config.SpotifyPlaylist),
		spottyChan:   spottyChan,
	}

	return spotifyHandler
}

// SpottyLoop is the main loop handling all spotify interactions
func (sh *SpotifyHandler) SpottyLoop() {
	// We will loop every time we get some IDs or the exit command from the spotty channel
	for {
		incomingMsg := <-sh.spottyChan

		switch incomingMsg.Kind {
		case types.AddedBy:
			songid := incomingMsg.ID
			addedBy, err := sh.db.getAddedBy(songid)
			response := types.SpottyMessage{
				Kind: types.AddedBy,
				User: addedBy,
				Err:  err,
			}
			sh.spottyChan <- response
			break
		case types.Track:
			songid := incomingMsg.ID
			user := incomingMsg.User

			if sh.db.isUnique(songid) {
				log.Println("Adding song to playlist")
				_, err := sh.spottyClient.AddTracksToPlaylist(sh.playlistID, spotify.ID(songid))
				if err != nil {
					log.Println(err)
					response := types.SpottyMessage{
						Err: err,
					}
					sh.spottyChan <- response
				} else {
					sh.db.addNewSong(songid, user)
					response := types.SpottyMessage{
						Msg: "I've added the song to the channel playlist!",
						Err: nil,
					}
					sh.spottyChan <- response
				}
			} else {
				response := types.SpottyMessage{
					Msg: "I can't add this song as it is already in the playlist according to my database.",
					Err: nil,
				}
				sh.spottyChan <- response
			}
			break
		}
	}
}

func newDatabase(config *types.Config) *database {
	db := &database{}
	db.init(config)
	db.setCollection(config)
	return db
}

func (db *database) init(config *types.Config) {
	// Create a mongodb client
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatalf("Error creating mongodb client: %v", err)
	}

	// Use context that times out after 10 seconds trying to connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to mongodb clinet
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Error connecting to mongodb client: %v", err)
	}

	// Verify our connection to mongodb client
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to mongodb")

	// Set our database client to the one we just made
	db.client = client
}

func (db *database) setCollection(config *types.Config) {
	db.collection = db.client.Database(config.MongoDatabase).Collection(config.MongoCollection)
}

func (db *database) addNewSong(songid string, user string) {
	// Use context that times out after 10 seconds trying to connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newSong := SongRecord{ID: songid, AddedBy: user}
	_, err := db.collection.InsertOne(ctx, newSong)

	if err != nil {
		log.Fatal(err)
	}
}

func (db *database) isUnique(songid string) bool {
	// Use context that times out after 10 seconds trying to connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result SongRecord
	searchTerm := bson.D{{
		Key:   "id",
		Value: songid,
	}}
	err := db.collection.FindOne(ctx, searchTerm).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This is the error we were hoping for, no documents exist of this song
			return true
		}
	}

	log.Println("Song already exists in database")

	return false
}

func (db *database) getAddedBy(songid string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result SongRecord
	searchTerm := bson.D{{
		Key:   "id",
		Value: songid,
	}}
	err := db.collection.FindOne(ctx, searchTerm).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "This song hasn't been added", nil
		}
		return "", err
	}

	return result.AddedBy, nil
}
