package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ResponseStruct []struct {
	Action     string   `json:"action"`
	Hostname   string   `json:"hostname"`
	Ipvers     int64    `json:"ipvers"`
	Namespace  string   `json:"namespace"`
	NexthopIps []string `json:"nexthopIps"`
	Oifs       []string `json:"oifs"`
	Preference int64    `json:"preference"`
	Prefix     string   `json:"prefix"`
	Protocol   string   `json:"protocol"`
	Source     string   `json:"source"`
	Timestamp  int64    `json:"timestamp"`
	Vrf        string   `json:"vrf"`
}

func main() {
	// Make the http request
	req, err := http.NewRequest("GET", "https://172.20.20.102:8000/api/v2/route/show?access_token=496157e6e869ef7f3d6ecb24a6f6d847b224ee4f", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Accept", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	//Read the response body
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Cannot marshall data", err)
	}
	// Create a pointer to the response struct
	JsonResponse := ResponseStruct{}
	//Unmarshall it
	err = json.Unmarshal(responseData, &JsonResponse)
	// Show how type values work.
	for _, v := range JsonResponse {
		fmt.Print("Prefix \t", v.Prefix+"\t", v.Protocol+"\t", "On Host \t", v.Hostname, "\n")
	}
}
