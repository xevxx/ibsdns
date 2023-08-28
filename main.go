package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// func check(e error) {
// 	if e != nil {
// 		fmt.Println(e)
// 	}
// }

func main() {

	value := flag.String("value", "", "IP to update domain to, example 192.168.1.1")
	flag.Parse()

	c, err := config()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}
	recordType := "A"
	newValue := *value

	if newValue == "" {
		fmt.Println("You must provide a value for the IP")
		os.Exit(1)
	}

	domainList := strings.Split(c.Domain, ",")
	for _, domain := range domainList {
		tid, status, message := updateDns(c.Url, c.ApiKey, c.Password, domain, recordType, newValue)
		if message != "" {
			fmt.Printf("Domain: %s, TransactID: %s, Status: %s, Message: %s\n", domain, tid, status, message)
		} else {
			fmt.Printf("Domain: %s, TransactID: %s, Status: %s\n", domain, tid, status)
		}
	}
}
