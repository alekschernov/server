package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Define the map data, where "G" is Grass and "W" is Water
var mapData = [][]string{
	{"G", "G", "G", "G", "G"},
	{"G", "W", "G", "W", "G"},
	{"G", "G", "P", "G", "G"}, // 'P' is the player position
	{"G", "G", "G", "G", "G"},
	{"G", "G", "G", "G", "G"},
}

// Response structure for the map
type MapResponse struct {
	Map [][]string `json:"map"`
}

// Request structure for pathfinding
type PathRequest struct {
	Start [2]int `json:"start"`
	Goal  [2]int `json:"goal"`
}

// Response structure for the calculated path
type PathResponse struct {
	Path [][2]int `json:"path"`
}

// A* pathfinding algorithm (simplified for example)
func astar(start [2]int, goal [2]int) [][2]int {
	var path [][2]int
	current := start
	path = append(path, current)

	for current != goal {
		// Check if we can move up, down, left, or right
		if current[0] < goal[0] && isWalkable(current[0]+1, current[1]) {
			current[0]++
		} else if current[0] > goal[0] && isWalkable(current[0]-1, current[1]) {
			current[0]--
		} else if current[1] < goal[1] && isWalkable(current[0], current[1]+1) {
			current[1]++
		} else if current[1] > goal[1] && isWalkable(current[0], current[1]-1) {
			current[1]--
		} else {
			// If there's no valid movement, break out of the loop to avoid infinite loop
			break
		}
		path = append(path, current)
	}
	return path
}

// Function to check if a tile is walkable
func isWalkable(row int, col int) bool {
	if row < 0 || col < 0 || row >= len(mapData) || col >= len(mapData[0]) {
		return false // Out of bounds
	}
	return mapData[row][col] != "W" // Check if the tile is not Water
}

func getMap(w http.ResponseWriter, r *http.Request) {
	resp := MapResponse{Map: mapData}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func calculatePath(w http.ResponseWriter, r *http.Request) {
	var req PathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate the path using the A* algorithm
	path := astar(req.Start, req.Goal)
	resp := PathResponse{Path: path}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/get-map", getMap)           // To serve map data
	http.HandleFunc("/calculate-path", calculatePath) // To calculate path

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}