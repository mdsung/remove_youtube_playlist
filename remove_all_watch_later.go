// Sample Go code for user authorization

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

// Load .env file
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

// func channelsListByUsername(service *youtube.Service, part []string, forUsername string) {
// 	call := service.Channels.List(part)
// 	call = call.ForUsername(forUsername)
// 	response, err := call.Do()
// 	handleError(err, "")

// 	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
// 		"and it has %d views.",
// 		response.Items[0].Id,
// 		response.Items[0].Snippet.Title,
// 		response.Items[0].Statistics.VideoCount))
// }

func removeItem(service *youtube.Service, id string) {
	call := service.PlaylistItems.Delete(id)
	err := call.Do()
	handleError(err, "")
}

func main() {

	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	playlistId := goDotEnvVariable("PLAYLIST_ID")
	fmt.Printf("Videos in list %s\r\n", playlistId)

	nextPageToken := ""

	for {
		// Call the playlistItems.list method to retrieve the
		// list of uploaded videos. Each request retrieves 50
		// videos until all videos have been retrieved.
		playlistCall := service.PlaylistItems.List([]string{"snippet"}).
			PlaylistId(playlistId).
			MaxResults(50).
			PageToken(nextPageToken)

		playlistResponse, err := playlistCall.Do()

		if err != nil {
			// The playlistItems.list method call returned an error.
			log.Fatalf("Error fetching playlist items: %v", err.Error())
		}

		for _, playlistItem := range playlistResponse.Items {
			title := playlistItem.Snippet.Title
			itemId := playlistItem.Id
			fmt.Printf("%v, (%v)\r\n", title, itemId)
			removeItem(service, itemId)
		}

		// Set the token to retrieve the next page of results
		// or exit the loop if all results have been retrieved.
		nextPageToken = playlistResponse.NextPageToken
		if nextPageToken == "" {
			break
		}
		fmt.Println()
	}
}
