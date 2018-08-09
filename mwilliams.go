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

// https://godoc.org/google.golang.org/api/sheets/v4#BatchUpdateValuesRequest
// https://developers.google.com/calendar/quickstart/go
// https://developers.google.com/sheets/quickstart/go
// https://console.developers.google.com/cloud-resource-manager

func getCreds() (*http.Client, *http.Client) {
	sheetsCred, err := ioutil.ReadFile("sheet_credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	calCred, err := ioutil.ReadFile("calendar_credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	sheetconfig, err := google.ConfigFromJSON(sheetsCred, "https://www.googleapis.com/auth/spreadsheets")
	calconfig, err := google.ConfigFromJSON(calCred, "https://www.googleapis.com/auth/calendar")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	sheetclient := getClient(sheetconfig)
	calclient := getClient(calconfig)

	if err != nil {
		log.Fatal("client error: ", err)
	}
	return sheetclient, calclient
}

func main() {

	ctx := context.Background()
	shclient, calclient := getCreds()

	sheetsService, err := sheets.New(shclient)
	calendarService, err := calendar.New(calclient)
	if err != nil {
		log.Fatal("services error: ", err)
	}
	// calendar id: Y3RhbGxhZGVuNzhAZ21haWwuY29t
	// event test: NnBscTZhOG01YzdqMnQ5NG0yYWQybmc4bzAgY3RhbGxhZGVuNzhAbQ
	// resp1, err := calendarService.Events.Get("ctalladen78", "M3RzNXUzMnV1N2dqN2I3c3FtaHZ1N3Y0YzcgY3RhbGxhZGVuNzhAbQ").Context(ctx).Do()
	resp1, err := calendarService.Events.Get("Y3RhbGxhZGVuNzhAZ21haWwuY29t", "NnBscTZhOG01YzdqMnQ5NG0yYWQybmc4bzAgY3RhbGxhZGVuNzhAbQ").Context(ctx).Do()

	if err != nil {
		log.Fatal("calendar error: ", err)
	}
	fmt.Println("resp1", calendarService)

	// TODO note calendar id to be used to get events ie: ctalladen78@gmail.com
	// TODO create a filter that returns only the details of the event/poeple on schedule within a time period
	// TODO create a new sheet on google drive, get sheet id
	// TODO create a double slice interface to populate the sheet
	reqBody := &sheets.Spreadsheet{
		// request body
		SpreadsheetId: "1A0IZHmHwjPf5YxLQAEQDuqsZj3rKW2VgDHZq9rmBalo",
	}
	resp, err := sheetsService.Spreadsheets.Create(reqBody).Context(ctx).Do()
	// resp, err := sheetsService.Spreadsheets.values.Get()
	// resp, err := sheetsService.Spreadsheets.values.update()
	spreadsheetID := "1A0IZHmHwjPf5YxLQAEQDuqsZj3rKW2VgDHZq9rmBalo"
	rangeData := "Sheet1!A1:B3"
	values := [][]interface{}{
		{"sample_A1", "sample_B1"},
		{"sample_A2", "sample_B2"},
		{"sample_A3", "sample_A3"}}
	// values := [][]interface{}{{"test"}}

	// https://godoc.org/google.golang.org/api/sheets/v4#BatchUpdateValuesRequest
	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}
	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  rangeData,
		Values: values,
	})
	_, err = sheetsService.Spreadsheets.Values.BatchUpdate(spreadsheetID, rb).Context(ctx).Do()
	// fmt.Println("complete", sheetsService.Spreadsheets)
	fmt.Println("complete: ", resp)
	fmt.Println("complete: ", resp1)
	// fmt.Println("complete: ", resp.SpreadsheetUrl)
}
