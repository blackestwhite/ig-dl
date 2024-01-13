package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
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

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}

	scriptContent := findScriptTag(doc)
	fmt.Println("Script Content:", scriptContent)
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

func findScriptTag(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "script" {
		// Check if the script tag contains the desired content (you may need to adjust this condition)
		if strings.Contains(n.FirstChild.Data, "your_desired_content_marker") {
			return n.FirstChild.Data
		}
	}

	// Recursively search for script tag
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findScriptTag(c); result != "" {
			return result
		}
	}

	return ""
}
