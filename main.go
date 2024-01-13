package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

const format = "https://instagram.com/p/%s/embed/captioned"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var input string
	fmt.Println("enter url:")
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatal(err)
	}

	postid, err := extractPostID(input)
	if err != nil {
		log.Fatal(err)
	}

	formatted := fmt.Sprintf(format, postid)

	client := &http.Client{}

	req, err := http.NewRequest("GET", formatted, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("HTTP request failed with status: %s", res.Status)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response Body:", string(body))
}

func extractPostID(url string) (string, error) {
	re := regexp.MustCompile(`https://www\.instagram\.com/reel/([a-zA-Z0-9_-]+)/`)

	matches := re.FindStringSubmatch(url)

	if len(matches) < 2 {
		return "", fmt.Errorf("post ID not found in the URL")
	}

	postID := matches[1]

	return postID, nil
}
