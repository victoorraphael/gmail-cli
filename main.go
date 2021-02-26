package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	// r, err := srv.Users.Labels.List(user).Do()
	// // r, err := srv.Users.Messages.List(user).Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve labels: %v", err)
	// }

	// for _, l := range r.Labels {
	// 	fmt.Println("Retorno de srv.Users", l.Name)
	// }

	m, errm := srv.Users.Messages.List(user).Do()

	if errm != nil {
		log.Fatalf("Unable to retrieve messages: %v", errm)
	}

	for _, ms := range m.Messages {
		msg, errmsg := srv.Users.Messages.Get(user, ms.Id).Do()

		if errmsg != nil {
			log.Fatalf("[#1] Unable to retrieve messages: %v", errm)
		}

		dec := new(mime.WordDecoder)
		data, err := dec.Decode(msg.Payload.Body.Data)

		if err != nil {
			log.Fatalf("[#2] Unable to retrieve messages: %v", errm)
		}

		fmt.Println("M ====>", data)
	}

	// if len(r.Messages) == 0 {
	// 	fmt.Println("No msgs found.")
	// 	return
	// }
	// fmt.Println("Messages:")
	// for _, l := range r.Messages {
	// 	fmt.Printf("- %s\n", l.Payload.MimeType)
	// }
}
