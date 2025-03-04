package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", helloWorldHandler)
	http.HandleFunc("/set", setUsernameHandler)
	http.HandleFunc("/get", getUsernameHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}

func setUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	key := r.URL.Query().Get("key")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	// Read the entire file into memory
	file, err := os.OpenFile("user.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "Unable to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if parts[0] == username {
			// Update the key if username exists
			lines = append(lines, username+","+key)
			found = true
		} else {
			lines = append(lines, line)
		}
	}

	if !found {
		// Append new username and key if not found
		lines = append(lines, username+","+key)
	}

	// Write the updated content back to the file
	file.Seek(0, 0)
	file.Truncate(0)
	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			http.Error(w, "Unable to write to file", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "saved successfully")
}

func getUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	file, err := os.Open("user.txt")
	if err != nil {
		http.Error(w, "Unable to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if parts[0] == username {
			fmt.Fprintf(w, "%s", parts[1])
			return
		}
	}
	fmt.Fprintf(w, "Username %s not found", username)

}
