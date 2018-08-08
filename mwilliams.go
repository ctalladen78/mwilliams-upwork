package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	context "golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	sheets "google.golang.org/api/sheets/v4"
)

func getCreds() *http.Client {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	if err != nil {
		log.Fatal("client error: ", err)
	}
	return client
}

func main() {

	ctx := context.Background()
	client := getCreds()

	sheetsService, err := sheets.New(client)
	calendarService, err := calendar.New(client)
	fmt.Println(calendarService)

	if err != nil {
		log.Fatal("sheets error: ", err)
	}

	// TODO note calendar id to be used to get events ie: ctalladen78@gmail.com
	// TODO create a filter that returns only the details of the event/poeple on schedule within a time period
	// TODO create a new sheet on google drive, get sheet id
	// TODO create a double slice interface to populate the sheet
	reqBody := &sheets.Spreadsheet{
		// request body
	}
	resp, err := sheetsService.Spreadsheets.Create(reqBody).Context(ctx).Do()

	fmt.Println("complete: %#v \n", resp)
}
