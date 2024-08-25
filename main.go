package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type SadhakaRequest struct {
	City    string `json:"city"`
	Lat     string `json:"lat"`
	Long    string `json:"long"`
	Sadhaka string `json:"sadhaka"`
}

type ParsedRequest struct {
	City    string
	Lat     float64
	Long    float64
	Sadhaka string
	Time    time.Time
}

type SadhakaResponse struct {
	NumberOfSadhakas int `json:"numberOfSadhakas"`
}

type SadhakaCounts struct {
	cities   map[string]int
	sadhakas map[string]time.Time
}

var counts SadhakaCounts = SadhakaCounts{
	cities:   make(map[string]int),
	sadhakas: make(map[string]time.Time),
}

func main() {
	fmt.Println("hello world")

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/sadhaka", serveSadhaka)
	http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	indexPath := filepath.Join("static", "index.html")
	http.ServeFile(w, r, indexPath)
}

func serveSadhaka(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req SadhakaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert Lat and Long from strings to float64
	lat, err := strconv.ParseFloat(req.Lat, 64)
	if err != nil {
		http.Error(w, "Invalid latitude format", http.StatusBadRequest)
		return
	}

	long, err := strconv.ParseFloat(req.Long, 64)
	if err != nil {
		http.Error(w, "Invalid longitude format", http.StatusBadRequest)
		return
	}

	// Process the data (this is just a placeholder, as no processing logic is provided)
	log.Printf("Received Sadhaka data: City=%s, Lat=%f, Long=%f, Sadhaka=%s", req.City, lat, long, req.Sadhaka)

	if err := appendToFile("sadhaka_requests.log", req); err != nil {
		http.Error(w, "Failed to write to log file", http.StatusInternalServerError)
		return
	}

	var parsedRequest = ParsedRequest{
		City:    req.City,
		Sadhaka: req.Sadhaka,
		Lat:     lat,
		Long:    long,
		Time:    time.Now(),
	}

	// Prepare the response
	sadhakaCount := updateAndGetSadhakaCount(parsedRequest, &counts)
	response := SadhakaResponse{
		NumberOfSadhakas: sadhakaCount,
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func updateAndGetSadhakaCount(req ParsedRequest, counts *SadhakaCounts) int {
	if counts.sadhakas[req.Sadhaka].IsZero() {
		counts.sadhakas[req.Sadhaka] = req.Time
	} else {
		return counts.cities[req.City]
	}
	counts.cities[req.City]++

	return counts.cities[req.City]
}

func appendToFile(filename string, req SadhakaRequest) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	logEntry := time.Now().Format(time.RFC3339) + " - " +
		"City: " + req.City + ", " +
		"Lat: " + req.Lat + ", " +
		"Long: " + req.Long + ", " +
		"Sadhaka: " + req.Sadhaka + "\n"

	if _, err := file.WriteString(logEntry); err != nil {
		return err
	}

	return nil
}
