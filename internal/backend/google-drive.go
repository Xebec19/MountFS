package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveRemote struct {
	root string // Google Drive folder ID to scope operations
	mu   sync.RWMutex
	srv  *drive.Service
}

// Example method using root to list files in the root folder
func (g *GoogleDriveRemote) ListFilesInRoot() ([]*drive.File, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	query := fmt.Sprintf("'%s' in parents and trashed = false", g.root)
	files := []*drive.File{}
	pageToken := ""
	for {
		req := g.srv.Files.List().Q(query).Fields("nextPageToken, files(id, name)").PageToken(pageToken)
		res, err := req.Do()
		if err != nil {
			return nil, err
		}
		files = append(files, res.Files...)
		if res.NextPageToken == "" {
			break
		}
		pageToken = res.NextPageToken
	}
	return files, nil
}

func NewGoogleDriveRemote(root string, credentialsPath string) (*GoogleDriveRemote, error) {
	ctx := context.Background()

	if credentialsPath == "" {
		credentialsPath = os.Getenv("GOOGLE_DRIVE_CREDENTIALS")
		if credentialsPath == "" {
			credentialsPath = "credentials.json"
		}
	}

	b, err := os.ReadFile(credentialsPath) // Downloaded from Google Cloud
	if err != nil {
		return nil, fmt.Errorf("Unable to read credentials file (%s): %v", credentialsPath, err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse config: %v", err)
	}

	client := getClient(config)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("Unable to create Drive client: %v", err)
	}

	return &GoogleDriveRemote{
		root: root,
		srv:  srv,
	}, nil
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Open this URL in browser:", authURL)

	var authCode string
	fmt.Print("Enter the code: ")
	fmt.Scan(&authCode)

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("Unable to create token file: %v", err)
		return
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		log.Printf("Unable to encode token to file: %v", err)
		return
	}
	fmt.Println("Saved token to", path)
}
