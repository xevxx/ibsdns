package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var mutex sync.Mutex

// getClientIP helps extract the client IP address from the request headers
func getClientIP(r *http.Request) string {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	}
	return ip
}

func loadLastIP() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	configPath := filepath.Join(configDir, "ibsdns", "lastIP.txt")
	var lastIP string
	// Check if file exists
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return ""
		}
		lastIP = strings.TrimSpace(string(data))
	}

	return lastIP
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
	// Extract the client's IP from the request
	//clientIP := strings.Split(r.RemoteAddr, ":")[0]
	lastIP := loadLastIP()

	clientIP := getClientIP(r)

	c, err := config()
	if err != nil {
		http.Error(w, "Unable to load config", http.StatusInternalServerError)
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
	http.HandleFunc("/update-dns", updateDnsHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
