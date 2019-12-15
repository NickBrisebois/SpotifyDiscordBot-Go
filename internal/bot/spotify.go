package bot

import (
	"context"
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/config"
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
func InitSpotify(config *config.Config, spottyChan chan string, client *spotify.Client) (err error) {
	var playlistID spotify.ID
	playlistID = spotify.ID(config.SpotifyPlaylist)

	if err != nil {
		log.Fatalf("couldn't get features playlists:%v", err)
	}

	songdb := newDatabase(config)

	// We will loop every time we get some IDs or the exit command from the spotty channel
	for {
		data := <-spottyChan

		if data == "quit" {
			break
		} else {
			if songdb.isUnique(data) {
				log.Println("Adding song to playlist")
				_, err := client.AddTracksToPlaylist(playlistID, spotify.ID(data))
				if err != nil {
					log.Println(err)
				} else {
					spottyChan <- "I've added the song to the channel playlist!"
				}
			} else {
				spottyChan <- "I can't add this song as it is already in the playlist according to my database."
			}
		}
	}

	return nil
}

func newDatabase(config *config.Config) *database {
	db := &database{}
	db.init(config)
	db.setCollection(config)
	return db
}

func (db *database) init(config *config.Config) {
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

func (db *database) setCollection(config *config.Config) {
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
			return true
		}
	}

	log.Println("Song already exists in database")
	return false
}
