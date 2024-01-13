package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const format = "https://instagram.com/p/%s/embed/captioned"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var input string
	fmt.Println("enter url:")
	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Println(err)
		return
	}

	postid, err := extractPostID(input)
	if err != nil {
		log.Println(err)
		return
	}

	formatted := fmt.Sprintf(format, postid)

	client := &http.Client{}

	req, err := http.NewRequest("GET", formatted, nil)
	if err != nil {
		log.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("HTTP request failed with status: %s", res.Status)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	doc := string(body)

	scriptContent := findLastScriptTagInBody(doc)

	videoURL, err := extractVideoURL(scriptContent)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Video URL:", videoURL)
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

func findLastScriptTagInBody(body string) string {
	// Use a regex to find the last script tag in the body
	re := regexp.MustCompile(`<script[^>]*>(.*?)</script>`)
	matches := re.FindAllStringSubmatch(body, -1)

	if len(matches) > 0 {
		return matches[len(matches)-1][1]
	}

	return ""
}

func extractVideoURL(scriptContent string) (string, error) {
	re := regexp.MustCompile(`\\"video_url\\":\\"(.*?)\\",`)
	matches := re.FindStringSubmatch(scriptContent)

	if len(matches) >= 2 {
		// Replace escaped slashes with regular slashes
		url := strings.Replace(matches[1], `\\\/`, `/`, -1)
		return url, nil
	}

	return "", fmt.Errorf("video_url not found in the script content")
}
