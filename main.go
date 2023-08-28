package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var mutex sync.Mutex

// getClientIP helps extract the client IP address from the request headers
func getClientIP(r *http.Request) (string, error) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			ip = strings.TrimSpace(ips[0]) // Using the first IP in the list
		}
	}

	if isValidIP(ip) {
		return ip, nil
	} else {
		return "", errors.New("invalid ip address")
	}
}

// isValidIP checks if the given string is a valid IPv4 or IPv6 address
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func loadLastIP() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(configDir, "ibsdns", "lastIP.txt")
	var lastIP string

	// Check if file exists
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return "", err
		}
		lastIP = strings.TrimSpace(string(data))
	} else if os.IsNotExist(err) {
		// Handle the case where the file does not exist
		return "", nil
	} else {
		// Handle other types of errors
		return "", err
	}

	return lastIP, nil
}

func saveLastIP(lastIP string) error {
	mutex.Lock()
	defer mutex.Unlock()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "ibsdns", "lastIP.txt")

	// Create the directory if it does not exist
	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}

	// Write the last IP to the file, creating the file if it doesn't exist
	return os.WriteFile(configPath, []byte(lastIP), 0644)
}

func updateDnsHandler(w http.ResponseWriter, r *http.Request) {

	apiKey := r.Header.Get("X-API-Key")

	clientIP, err := getClientIP(r)
	if err != nil {
		http.Error(w, "Unable to get client IP", http.StatusInternalServerError)
		return
	}

	// Load last known IP
	lastIP, err := loadLastIP()
	if err != nil && lastIP != "" {
		http.Error(w, "Unable to load last known IP", http.StatusInternalServerError)
		return
	}

	c, err := config()
	if err != nil {
		http.Error(w, "Unable to load config", http.StatusInternalServerError)
		return
	}

	if apiKey != c.ApiKey { // Pass the expected API key from config
		http.Error(w, "Invalid API Key", http.StatusUnauthorized)
		return
	}

	// Dummy variables, you can replace these with real values
	recordType := "A"

	domainList := strings.Split(c.Domain, ",")
	if lastIP == clientIP {
		fmt.Fprintf(w, "Domain: %s, No change in IP address. No update needed.\n", clientIP)
	} else {
		lastIP = clientIP // Update the last known IP
		saveLastIP(lastIP)
		for _, domain := range domainList {
			tid, status, message := updateDns(c.Url, c.ApiKeyInternetBS, c.Password, domain, recordType, clientIP)

			// Dummy values for tid, status, message to fake the DNS update
			//tid, status, message := "12345", "Success", ""
			//fmt.Fprintf(w, "RecordType: %s, record type.\n", recordType)

			if message != "" {
				fmt.Fprintf(w, "Domain: %s, TransactID: %s, Status: %s, Message: %s\n", domain, tid, status, message)
			} else {
				fmt.Fprintf(w, "Domain: %s, TransactID: %s, Status: %s\n", domain, tid, status)
			}
		}
	}
}

func main() {
	c, err := config()
	if err != nil {
		return
	}
	host := c.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := c.Port
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/update-dns", updateDnsHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(host+":"+port, nil)
}
