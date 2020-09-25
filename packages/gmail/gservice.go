package gservice

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"database/sql"
	"meli/domain/sql"
)

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		tok = GetTokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
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

// Retrieves a token from a local file at project
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path of project
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func CreateService() (*gmail.Service, error) {
	credentialsPath, present := os.LookupEnv("CREDENTIALS_JSON_GMAIL")
	if !present {
		log.Fatalf("CREDENTIALS_JSON_GMAIL env var is mandatory")
	}
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := GetClient(config)

	return gmail.New(client)
}

func FindMessages(query string, srv *gmail.Service, db *sql.DB) {
	token := ""
	for {
		rr, err := srv.Users.Messages.List("me").PageToken(token).Q(query).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve messages: %v", err)
		}
		if len(rr.Messages) == 0 {
			fmt.Println("No messages found.")
			return
		}

		for _, m := range rr.Messages {
			msg, errM := srv.Users.Messages.Get("me", m.Id).Format("metadata").MetadataHeaders("Subject", "Date", "From").Do()
			if errM != nil {
				log.Fatalf("Unable to retrieve messages: %v", err)
			}

			for _, header := range msg.Payload.Headers {
				if strings.Contains(header.Name, "Subject") {
					showMessageFromPayload(string(header.Value))
					chooseIfPersist(msg, db)
					break
				}
			}
		}

		token = rr.NextPageToken

		if token == "" {
			break
		}
	}
}
func showMessageFromPayload(msg string) {
	fmt.Printf("Message found with subject: %s,\n", msg)
}
func chooseIfPersist(msg *gmail.Message, db *sql.DB) {
	fmt.Print("Do you want to persist the message? (y/n/q): ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	if strings.TrimRight(text, "\r\n") == "y" || strings.TrimRight(text, "\r\n") == "Y" {

		fecha, from_, subject, errFields := extractFields(msg.Payload.Headers)
		if errFields != nil {
			panic(errFields)
		}
		persistedId := dao.Persist(db, fecha, from_, subject)
		fmt.Printf("Persisted: %s\n", persistedId)
		return
	}
	if strings.TrimRight(text, "\r\n") == "n" || strings.TrimRight(text, "\r\n") == "N" {
		return
	}
	if strings.TrimRight(text, "\r\n") == "q" || strings.TrimRight(text, "\r\n") == "Q" {
		os.Exit(0)
	}
}

func extractFields(headers []*gmail.MessagePartHeader) (string, string, string, error) {
	c := 0
	var fecha, from_, subject string
	var error_ error
	for _, header := range headers {
		if c == 3 {
			break
		}
		if header.Name == "Date" {
			c = c + 1
			fecha = header.Value
			continue
		}
		if header.Name == "From" {
			c = c + 1
			from_ = header.Value
			continue
		}
		if header.Name == "Subject" {
			c = c + 1
			subject = header.Value
			continue
		}
		if c == 3 {
			break
		}
	}
	if c != 3 {
		error_ = fmt.Errorf("Fields fecha, from or subject not found")
	}
	return fecha, from_, subject, error_
}

func showMessage(msg *gmail.Message) {
	out, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Message: %s\n", string(out))
}
