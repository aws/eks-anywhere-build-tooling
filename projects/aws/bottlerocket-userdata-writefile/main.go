package main

import (
	"io/ioutil"
	"log"
	"net/http"
	url "net/url"
	"os"
	"path"
)

const hegelUserDataVersion = "2009-04-04"

func main() {
	log.Print("Bottlerocket UserData - Starting process")

	hegelUrl := os.Getenv("HEGEL_URL")
	if hegelUrl == "" {
		log.Fatalf("Error - No HEGEL_URL env var found")
	}

	url, err := url.Parse(hegelUrl)
	if err != nil {
		log.Fatalf("Error parsing hegel url: %v", err)
	}
	url.Path = path.Join(url.Path, hegelUserDataVersion, "user-data")
	resp, err := http.Get(url.String())
	if err != nil {
		log.Fatalf("Error with HTTP GET call: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading HTTP GET response body: %v", err)
	}

	log.Println(string(respBody))
}
