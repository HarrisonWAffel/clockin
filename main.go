package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	//"strings"
	"time"
)

func main() {
	//this program uses the Google sheets API to fill
	//a spreadsheet for my work hours

	//req's:
	// 	gotta write start and end times
	//	gotta be able to read the difference between the two and fill that out as well


	//~~ setup the api ~~//
	if len(os.Args) == 1 {
		fmt.Println("Please Supply Required Parameters")
		fmt.Println("s : setup, link to a new spread sheet or account")
		fmt.Println("i : clock in, write starting time to sheet")
		fmt.Println("o : clock out, write ending time to sheet")
		return
	}

		b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}



	if os.Args[1] == "s" {
		//if we want to setup
		f, err := os.Create("config.txt")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Please Enter the sheets ID you want to use ")

		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println("here")
		fmt.Println(text)
		_, err2 := f.WriteString(strings.TrimSuffix(text, "\n"))

		if err2 != nil {
			fmt.Println(err)
			_ = f.Close()
			return
		}
	}


	file, err := os.Open("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	t, err2 := ioutil.ReadAll(file)
	if err2 != nil{

	}

	defer file.Close()


	spreadsheetId := string(t)
	clockInRange := "Sheet1!A2:B99"
	clockOutRange := "Sheet1!D2:E99"
	formulaRange := "Sheet1!F2:G99"
	now := time.Now()

	//currentTime := now.Format("15:05:03")

	if os.Args[1] == "i" {
		fmt.Println("")

		var vr sheets.ValueRange
		vr.MajorDimension = "ROWS"
		vr.Range = clockInRange
		myval := []interface{}{now.Format("01-02-2006"), now.Format("03:04:05")}
		vr.Values = append(vr.Values, myval)
		_, err = srv.Spreadsheets.Values.Append(spreadsheetId, clockInRange, &vr).ValueInputOption("RAW").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet. %v", err)
		}else{
			fmt.Println(now.Format("01-02-2006") +" ] Succesfully Clocked in @ " + now.Format("03:04:05"))
		}
	}else


	if os.Args[1] == "o" {
		//step1 read the starting time to calc diff

		resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, clockInRange).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet: %v", err)

		}

		resp2, err2 := srv.Spreadsheets.Values.Get(spreadsheetId, clockOutRange).Do()
		if err2 != nil{
			log.Fatalf("Unable to retrieve data from sheet: %v", err)
		}

		if len(resp.Values) == 0 {
			fmt.Println("No data found.")
		} else {

			currentDate := now.Format("01-02-2006")
			startDate := resp.Values[len(resp.Values)-1][0]



			if startDate == currentDate && len(resp.Values)-1 == len(resp2.Values){

			//NOTE: the below code was an attempt to calc the difference in time within golang, but its easier to just let google sheets do this for us on the ss not the script
			//the code will stay for reference
			//	//we should only fill the value for the current day
			//	//this should also stop us from writing into the ClockInRange
			//
			//	//POSSIBLE ISSUES WITH DESIGN:
			//	//	What if we forget to clock out? the program would never fill in that gap, or even worse fill in the gap with bad data
			//	// this is an edge case, but will probably happen at some point
			//
			//	fullString := fmt.Sprintf("%v", startDate) + "T" + fmt.Sprintf("%v", startTime)
			//	fmt.Println(fullString)
			//	//take a string and split it, turning each piece into an int
			//	//t := fmt.Sprintf("%d", strings.Split(fmt.Sprintf("%v", startDate), "-"))
			//	//s := fmt.Sprintf("%d", strings.Split(fmt.Sprintf("%v", startTime), ":"))
			//
			//
			//	fmt.Println(now.Format("03:04:05"))
			//
			//	stTime, _ := time.Parse("03:04:05", fmt.Sprintf("%v", startTime))
			//	//stTime  := time.Date(int(t[2]), int(t[0]), int(t[1]))
			//
			//
			//
			//
			//	diff := stTime.Sub(now)
			//	fmt.Println(diff)

			var vr sheets.ValueRange
			vr.MajorDimension = "COLUMNS"
			myval := []interface{}{now.Format("03:04:05")}
			vr.Values = append(vr.Values, myval)
			vr.Range = clockOutRange
			_, err = srv.Spreadsheets.Values.Append(spreadsheetId, clockOutRange, &vr).ValueInputOption("RAW").Do()

			if err != nil {
				log.Fatalf("Unable to Write end time %v", err)
			}

			var ve sheets.ValueRange
			ve.MajorDimension = "COLUMNS"
			myval = []interface{}{"=D2:D10001-B2:B10001"}
			ve.Range = formulaRange
			ve.Values = append(ve.Values, myval)
			_, err = srv.Spreadsheets.Values.Append(spreadsheetId, formulaRange, &ve).ValueInputOption("USER_ENTERED").Do()

			if err != nil {
				log.Fatalf("Unable to write formula. %v", err)
			}else{
				fmt.Println(now.Format("01-02-2006") +" ] Succesfully Clocked out @ " + now.Format("03:04:05"))
			}


			}else{
				fmt.Println("You must clock in before you can clock out!  ")
				fmt.Println("If you Forgot to clock out in days prior YOU MUST FIX THIS BEFORE CONTINUING USE")
				fmt.Println("Hours Link: PUT HOURS LINK HERE")
			}
		}
	}else{
		fmt.Println("Please Supply A valid parameter.")
		fmt.Println("s : setup, link to a new spread sheet or account")
		fmt.Println("i : clock in, write starting time to sheet")
		fmt.Println("o : clock out, write ending time to sheet")
	}


}


//~~ core sheets api funcs ~~//
//read more https://developers.google.com/sheets/api/quickstart/go


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