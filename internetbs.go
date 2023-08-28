package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func updateDns(urlConfig, apiKey, password, fullrecordname, recordtype, newvalue string) (transactid, status, message string) {

	// Test Api
	Url := urlConfig

	param := url.Values{}
	param.Add("apiKey", apiKey)
	param.Add("password", password)
	param.Add("fullrecordname", fullrecordname)
	param.Add("type", recordtype)
	param.Add("newvalue", newvalue)
	param.Add("ResponseFormat", "JSON")

	req, err := http.NewRequest("POST", Url, strings.NewReader(param.Encode()))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	jsonblob, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	transactid, status, message = decodeResponse(jsonblob)
	return
}

func decodeResponse(jsonblob []byte) (transactid, status, message string) {
	var response Response
	json.Unmarshal(jsonblob, &response)
	transactid = response.Transactid
	status = response.Status
	message = response.Message

	return
}
