package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var successfulJokes []JokeResponse

type JokeResponse struct {
	Error    bool   `json:"error"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
	Flags    struct {
		NSFW      bool `json:"nsfw"`
		Religious bool `json:"religious"`
		Political bool `json:"political"`
		Racist    bool `json:"racist"`
		Sexist    bool `json:"sexist"`
		Explicit  bool `json:"explicit"`
	} `json:"flags"`
	ID       int      `json:"id"`
	Safe     bool     `json:"safe"`
	Lang     string   `json:"lang"`
	CausedBy []string `json:"causedBy"`
}

// first function to run
func main() {
	menu()
}

// concatenateStrings takes a slice of strings and concatenates them into a single string separated by commas which is later used to build a joke with multiple categories
func concatenateStrings(stringArray []string) string {

	result := strings.Join(stringArray, ",")

	return result
}

// buildAJokeWithCategory takes a slice of strings and builds a joke with multiple categories
func buildAJokeWithCategory(category []string) {
	if len(category) == 0 {
		category = append(category, "Any")
	}

	concatenatedString := concatenateStrings(category)

	resp, err := http.Get(getURL() + "/" + concatenatedString + "?blacklistFlags=nsfw,religious,political,racist,sexist,explicit")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var joke JokeResponse
	err = json.Unmarshal(body, &joke)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if joke.Error {
		fmt.Println("API error:", joke.CausedBy[0])
		return
	}

	successfulJokes = append(successfulJokes, joke)

	fmt.Println("Setup:", joke.Setup)
	fmt.Println("Delivery:", joke.Delivery)
}

// getURL returns the URL for the joke API
func getURL() string {
	return "https://v2.jokeapi.dev/joke"
}

// categoryMenu displays the category menu and allows the user to select a category to build a joke with
func categoryMenu() {

	var userInput string
	var category []string

	for userInput != "6" {
		fmt.Println("\nCategory Menu")
		fmt.Println("-------------------------")
		fmt.Println("1. Programming")
		fmt.Println("2. Miscellaneous")
		fmt.Println("3. Dark")
		fmt.Println("4. Pun")
		fmt.Println("5. Holiday")
		fmt.Println("6. Build Joke with Multiple Categories")
		fmt.Println("7. Back to Main Menu")

		fmt.Print("Enter your choice: ")
		fmt.Scanln(&userInput)
		fmt.Println()

		switch userInput {
		case "1":
			fmt.Println("Adding Programming to category list")
			category = append(category, "Programming")
		case "2":
			fmt.Println("Adding Miscellaneous to category list")
			category = append(category, "Miscellaneous")
		case "3":
			fmt.Println("Adding Dark to category list")
			category = append(category, "Dark")
		case "4":
			fmt.Println("Adding Pun to category list")
			category = append(category, "Pun")
		case "5":
			fmt.Println("Adding Holiday to category list")
			category = append(category, "Holiday")

		case "6":
			fmt.Println("Building joke with multiple categories")
			buildAJokeWithCategory(category)
		case "7":
			fmt.Println("Returning to Main Menu...")
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 6.")
		}
	}

}

// menu displays the main menu and allows the user to select an option
func menu() {
	var userInput string

	for userInput != "8" {
		fmt.Println("\nMain Menu")
		fmt.Println("-------------------------")
		fmt.Println("1. Get Any Joke")
		fmt.Println("2. Build a Joke with a Specific Word")
		fmt.Println("3. Build a Joke with a Category")
		fmt.Println("4. Save Jokes to File")
		fmt.Println("5. See all Successful Jokes")
		fmt.Println("6. Read Jokes from File")
		fmt.Println("7. Delete all Jokes from File")
		fmt.Println("8. Exit")
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&userInput)
		fmt.Println()

		switch userInput {
		case "1":
			getAnyJoke()
		case "2":
			fmt.Println("Building a joke with a specific word")
			fmt.Print("Enter a word: ")
			var word string
			fmt.Scanln(&word)
			getJokeWithWord(word)
		case "3":
			categoryMenu()
		case "4":
			fmt.Println("Saving jokes to file...")
			saveJokesToFile("jokes.json")
		case "5":
			fmt.Println("All Successful Jokes...")
			displayAllJokes()
		case "6":
			fmt.Println("Reading jokes from file...")
			readJokesFromFile("jokes.json")
		case "7":
			fmt.Println("Deleting all jokes from file...")
			deleteJokesFromFile("jokes.json")
		case "8":
			fmt.Println("Exiting...")
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 8.")
		}
	}
}

// getJokeWithWord gets a joke with a specific word
func getJokeWithWord(word string) {
	resp, err := http.Get(getURL() + "/Any?blacklistFlags=nsfw,religious,political,racist,sexist,explicit&contains=" + word)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var joke JokeResponse
	err = json.Unmarshal(body, &joke)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if joke.Error {
		fmt.Println("API error:", joke.CausedBy[0])
		return
	}
	if joke.Setup == "" {
		fmt.Println("No jokes found with the word:", word)
		return
	}

	successfulJokes = append(successfulJokes, joke)

	fmt.Println("Setup:", joke.Setup)
	fmt.Println("Delivery:", joke.Delivery)
}

// getAnyJoke gets a joke with any category
func getAnyJoke() {
	resp, err := http.Get(getURL() + "/Any?blacklistFlags=nsfw,religious,political,racist,sexist,explicit")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var joke JokeResponse
	err = json.Unmarshal(body, &joke)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if joke.Error {
		fmt.Println("API error:", joke.Setup)
		return
	}

	successfulJokes = append(successfulJokes, joke)

	fmt.Println("Setup:", joke.Setup)
	fmt.Println("Delivery:", joke.Delivery)
}

// displayAllJokes displays all jokes that have been successfully retrieved
func displayAllJokes() {
	fmt.Println("All Jokes:")
	fmt.Println("-------------------------")
	for i, joke := range successfulJokes {
		fmt.Printf("Joke %d:\n", i+1)
		fmt.Printf("Setup: %s\n", joke.Setup)
		fmt.Printf("Delivery: %s\n", joke.Delivery)
		fmt.Println("-------------------------")
	}
}

// saveJokesToFile saves all jokes that have been successfully retrieved to a file
func saveJokesToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	jokesJSON, err := json.MarshalIndent(successfulJokes, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling jokes to JSON:", err)
		return
	}

	_, err = file.Write(jokesJSON)
	if err != nil {
		fmt.Println("Error writing jokes to file:", err)
		return
	}

	fmt.Println("Jokes saved to", filename)
}

// readJokesFromFile reads jokes from a file and displays them
func readJokesFromFile(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var jokes []JokeResponse
	err = json.Unmarshal(data, &jokes)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("Jokes read from", filename)
	fmt.Println("-------------------------")
	for i, joke := range jokes {
		fmt.Printf("Joke %d:\n", i+1)
		fmt.Printf("Setup: %s\n", joke.Setup)
		fmt.Printf("Delivery: %s\n", joke.Delivery)
		fmt.Println("-------------------------")
	}
}

// deleteJokesFromFile deletes all jokes from a file
func deleteJokesFromFile(filename string) {

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fmt.Println("All jokes deleted from", filename)
}
