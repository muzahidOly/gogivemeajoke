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

func main() {
	menu()
}

func concatenateStrings(stringArray []string) string {
	// Join the strings with a comma in between
	result := strings.Join(stringArray, ",")

	return result
}

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

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse the JSON response
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

	// Print the setup and delivery
	fmt.Println("Setup:", joke.Setup)
	fmt.Println("Delivery:", joke.Delivery)
}

func getURL() string {
	return "https://v2.jokeapi.dev/joke"
}

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

func menu() {
	var userInput string

	for userInput != "7" {
		fmt.Println("\nMain Menu")
		fmt.Println("-------------------------")
		fmt.Println("1. Get Any Joke")
		fmt.Println("2. Build a Joke with a Specific Word")
		fmt.Println("3. Build a Joke with a Category")
		fmt.Println("4. Save Jokes to File")
		fmt.Println("5. See all Successful Jokes")
		fmt.Println("6. Read Jokes from File")
		fmt.Println("7. Exit")
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
			fmt.Println("Exiting...")
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 7.")
		}
	}
}

func getJokeWithWord(word string) {
	resp, err := http.Get(getURL() + "/Any?blacklistFlags=nsfw,religious,political,racist,sexist,explicit&contains=" + word)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse the JSON response
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

func getAnyJoke() {
	resp, err := http.Get(getURL() + "/Any?blacklistFlags=nsfw,religious,political,racist,sexist,explicit")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse the JSON response
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

func saveJokesToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Marshal the jokes to JSON format
	jokesJSON, err := json.MarshalIndent(successfulJokes, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling jokes to JSON:", err)
		return
	}

	// Write the JSON data to the file
	_, err = file.Write(jokesJSON)
	if err != nil {
		fmt.Println("Error writing jokes to file:", err)
		return
	}

	fmt.Println("Jokes saved to", filename)
}

func readJokesFromFile(filename string) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal the JSON data into a slice of JokeResponse structs
	var jokes []JokeResponse
	err = json.Unmarshal(data, &jokes)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Display the jokes
	fmt.Println("Jokes read from", filename)
	fmt.Println("-------------------------")
	for i, joke := range jokes {
		fmt.Printf("Joke %d:\n", i+1)
		fmt.Printf("Setup: %s\n", joke.Setup)
		fmt.Printf("Delivery: %s\n", joke.Delivery)
		fmt.Println("-------------------------")
	}
}
