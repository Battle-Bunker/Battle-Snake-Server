package main

import (
	"encoding/json"
	"fmt"
	"github.com/BattlesnakeOfficial/rules/cli/commands"
	"log"
	"net/http"
    "os"
)


// Amend PlayRequest to accept a collection of players
type PlayRequest struct {
	Players []commands.Player `json:"players"`
}

// Define a struct to hold the JSON data for index response
type IndexResponse struct {
	Status string `json:"status"`
}

func main() {
	// Set up a route and handler function
	http.HandleFunc("/play", playHandler)

	// Set up a route and handler function for the index page
	http.HandleFunc("/", indexHandler)

	// Start the HTTP server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// playHandler handles the POST request to the '/play' endpoint
func playHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var req PlayRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Create a temporary file for the game output
    tmpFile, err := os.CreateTemp("", "game_output_*.json")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating temp file: %v", err), http.StatusInternalServerError)
        return
    }
    defer os.Remove(tmpFile.Name()) // Clean up the file afterwards

    // Call PlayBattlesnakeGame with the collection of players and output path
    err = commands.PlayBattlesnakeGame(req.Players, tmpFile.Name())
    if err != nil {
        http.Error(w, fmt.Sprintf("Error playing game: %v", err), http.StatusInternalServerError)
        return
    }

    // Send the game file as a response
    http.ServeFile(w, r, tmpFile.Name())
}

// indexHandler handles the GET requests to the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	// Create response object
	response := IndexResponse{Status: "Ok"}
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	// Encode the response as JSON and send it
	json.NewEncoder(w).Encode(response)
}
