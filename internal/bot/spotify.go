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

// SongRecord is what is saved to mongodb to keep track of what songs have been added
type SongRecord struct {
	ID string
}

// InitSpotify starts spotify API handler
func InitSpotify(config *types.Config, spottyChan chan types.SpottyMessage, client *spotify.Client) (err error) {
	var playlistID spotify.ID
	playlistID = spotify.ID(config.SpotifyPlaylist)

	if err != nil {
		log.Fatalf("couldn't get features playlists:%v", err)
	}

	songdb := newDatabase(config)

	// We will loop every time we get some IDs or the exit command from the spotty channel
	for {
		incomingMsg := <-spottyChan
		songid := incomingMsg.ID

		if songid == "quit" {
			break
		} else {
			if songdb.isUnique(songid) {
				log.Println("Adding song to playlist")
				_, err := client.AddTracksToPlaylist(playlistID, spotify.ID(songid))
				if err != nil {
					log.Println(err)
					response := types.SpottyMessage{
						Err: err,
					}
					spottyChan <- response
				} else {
					songdb.addNewSong(songid)
					response := types.SpottyMessage{
						Msg: "I've added the song to the channel playlist!",
					}
					spottyChan <- response
				}
			} else {
				response := types.SpottyMessage{
					Msg: "I can't add this song as it is already in the playlist according to my database.",
				}
				spottyChan <- response
			}
		}
	}

	return nil
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

func (db *database) addNewSong(songid string) {
	// Use context that times out after 10 seconds trying to connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newSong := SongRecord{ID: string(songid)}
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
