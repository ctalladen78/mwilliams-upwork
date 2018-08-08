
package main

import (
	"encoding/json"
    "fmt"
    "io/ioutil"
    //"errors"
    "log"
    "net/http"
    "os"
    context "golang.org/x/net/context"
   "golang.org/x/oauth2" 
   google "golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

// get a token and config http client
func getClient(config *oauth2.Config) (*http.Client) {
	// TODO: Change placeholder below to get authentication credentials. See
        // https://developers.google.com/sheets/quickstart/go#step_3_set_up_the_sample
        //
        // Authorize using the following scopes:
        //     sheets.DriveScope
        //     sheets.DriveFileScope
        //     sheets.SpreadsheetsScope
        //return nil, errors.New("not implemented")
        tokFile := "token.json"
        tok, err := getTokenFromFile(tokFile)
        if err != nil {
        	tok = getTokenFromWeb(config) // if doesnt exist
        	saveToken(tokFile, tok)
        }
        return config.Client(context.Background(), tok)
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
        fmt.Printf("Saving credential file to: %s\n", path)
        f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
        defer f.Close()
        if err != nil {
                log.Fatalf("Unable to cache oauth token: %v", err)
        }
        json.NewEncoder(f).Encode(token)
}

// generate a token file after authorizing from web
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("Unable to read authorization code: %v", err)
        }

        tok, err := config.Exchange(oauth2.NoContext, authCode)
        if err != nil {
                log.Fatalf("Unable to retrieve token from web: %v", err)
        }
        return tok
}

func getTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
        defer f.Close()
        if err != nil {
                return nil, err
        }
        tok := &oauth2.Token{}
        err = json.NewDecoder(f).Decode(tok)
        return tok, err
}

func main() {
	ctx := context.Background()
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
	sheetsService, err := sheets.New(client)
	if err != nil {
		log.Fatal("sheets error: ", err)
	}
	reqBody := &sheets.Spreadsheet{
		// add fields on request body
		/*
		{
			"sheets":[
				{
					"data": [
						{
							"startColumn": 0,
							"startRow": 0,
							"rowData": [
								{
									"values": [
										"effectiveValue": {
											"stringValue": "hello hello"
										}
									]
								}
							]
						}
					]
				}
			]
		}
		
	}
	*/
	resp, err := sheetsService.Spreadsheets.Create(reqBody).Context(ctx).Do()
	
	fmt.Println("complete: %#v \n", resp)
}
