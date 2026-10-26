package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	listImgsUrl := "https://api.cloudflare.com/client/v4/accounts/400b3f859551ad3d06330ed21b0630dc/images/v2"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}


	req, err := http.NewRequest("GET", listImgsUrl, nil)
	req.Header.Add("Authorization", os.Getenv("imgs_bearer"))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed GET request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed with HTTP error code: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var imageResponse ImageResponse
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		log.Fatalf("Failed to unmarshal %v", err)
	}

	images := imageResponse.Result.Images
	fmt.Println(len(images))
}
