package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// fetchContactInfo retrieves contact email and a short summary from a website URL.
func fetchContactInfo(url string) (contact string, summary string, err error) {
	if url == "" {
		return "", "", fmt.Errorf("no website URL")
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch website: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read website: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// try meta description or first paragraph as summary
	summary, _ = doc.Find("meta[name='description']").Attr("content")
	if summary == "" {
		summary, _ = doc.Find("meta[property='og:description']").Attr("content")
	}
	if summary == "" {
		summary = strings.TrimSpace(doc.Find("p").First().Text())
	}

	// extract first email address
	re := regexp.MustCompile(`[A-Za-z0-9._%%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`)
	match := re.FindString(string(body))
	contact = match

	return contact, summary, nil
}
