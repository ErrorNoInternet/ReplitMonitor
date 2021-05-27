package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type replJSON struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Language    string `json:"language"`
	TimeCreated string `json:"time_created"`
	TimeUpdated string `json:"time_updated"`
	Url         string `json:"url"`
}

type profileJSON struct {
	EmailHash     string     `json:"emailHash"`
	Origanization string     `json:"origanization"`
	Repls         []replJSON `json:"repls"`
	ID            int        `json:"id"`
	Username      string     `json:"username"`
	FirstName     string     `json:"firstName"`
	LastName      string     `json:"lastName"`
	Bio           string     `json:"bio"`
	IsTeam        bool       `json:"isTeam"`
	TopLanguages  []string   `json:"topLanguages"`
}

var expectedRepls []replJSON
var expectedUsername string
var expectedFirstName string
var expectedLastName string
var expectedBio string
var firstRun bool = true
var firstCheck bool = true

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Replit Username: ")
	scanner.Scan()
	userName := scanner.Text()

	for {
		fmt.Printf("\r[%v] Updating profile information...", time.Now().Format("2006-01-02 15:04:05"))
		var found bool = false
		response, errorObject := http.Get("https://replit.com/data/profiles/" + userName)
		if errorObject != nil {
			fmt.Println("Error: " + errorObject.Error())
		}
		dataBytes, errorObject := ioutil.ReadAll(response.Body)
		if errorObject != nil {
			fmt.Println("Error: " + errorObject.Error())
		}
		var profile profileJSON
		json.Unmarshal(dataBytes, &profile)

		if firstCheck {
			expectedBio = profile.Bio
			expectedFirstName = profile.FirstName
			expectedLastName = profile.LastName
			expectedUsername = profile.Username
			firstCheck = false
		}

		if firstRun {
			expectedRepls = append(expectedRepls, profile.Repls...)
			firstRun = false
		}

		if expectedUsername != profile.Username {
			fmt.Println("")
			fmt.Printf("\n[%v] User has changed username:\nBefore: %v\nAfter: %v\n\n", time.Now().Format("2006-01-02 15:04:05"), expectedUsername, profile.Username)
			expectedUsername = profile.Username
		}

		if expectedFirstName != profile.FirstName {
			fmt.Println("")
			fmt.Printf("\n[%v] User has changed firstName:\nBefore: %v\nAfter: %v\n\n", time.Now().Format("2006-01-02 15:04:05"), expectedFirstName, profile.FirstName)
			expectedFirstName = profile.FirstName
		}

		if expectedLastName != profile.LastName {
			fmt.Println("")
			fmt.Printf("\n[%v] User has changed lastName:\nBefore: %v\nAfter: %v\n\n", time.Now().Format("2006-01-02 15:04:05"), expectedLastName, profile.LastName)
			expectedLastName = profile.LastName
		}

		if expectedBio != profile.Bio {
			fmt.Println("")
			fmt.Printf("\n[%v] User has changed bio:\nBefore: %v\nAfter: %v\n\n", time.Now().Format("2006-01-02 15:04:05"), expectedBio, profile.Bio)
			expectedBio = profile.Bio
		}

		for _, repl := range profile.Repls {
			for _, expected := range expectedRepls {
				if repl.ID == expected.ID {
					found = true
					if repl.TimeUpdated != expected.TimeUpdated {
						fmt.Println("")
						fmt.Printf("\n[%v] Repl (%v) timeUpdated has been changed:\nBefore: %v\nAfter: %v\n\n", time.Now().Format("2006-01-02 15:04:05"), repl.Title, expected.TimeUpdated, repl.TimeUpdated)
						expectedRepls = []replJSON{}
						firstRun = true
					}
					if repl.Description != expected.Description {
						fmt.Println("")
						fmt.Printf("\n[%v] Repl (%v) description has been changed:\nBefore: %v\nAfter: %v\n\n", time.Now().Format("2006-01-02 15:04:05"), repl.Title, expected.Description, repl.Description)
						expectedRepls = []replJSON{}
						firstRun = true
					}
				}
			}
			if !found {
				expectedRepls = append(expectedRepls, repl)
				fmt.Println("")
				fmt.Printf("\n[%v] User created new REPL:\nName: %v\nTime: %v\nLanguage: %v\n\n", time.Now().Format("2006-01-02 15:04:05"), repl.Title, repl.TimeCreated, repl.Language)
			}
		}
		time.Sleep(20 * time.Second)
	}
}
