package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FetchBatchData fetches startup data for a given batch name from the YC API.
func FetchBatchData(batchName string) (startups []Startup, err error) { // Named return error 'err'
	const (
		appID  = "45BWZJ1SGC"
		apiKey = "MjBjYjRiMzY0NzdhZWY0NjExY2NhZjYxMGIxYjc2MTAwNWFkNTkwNTc4NjgxYjU0YzFhYTY2ZGQ5OGY5NDMxZnJlc3RyaWN0SW5kaWNlcz0lNUIlMjJZQ0NvbXBhbnlfcHJvZHVjdGlvbiUyMiUyQyUyMllDQ29tcGFueV9CeV9MYXVuY2hfRGF0ZV9wcm9kdWN0aW9uJTIyJTVEJnRhZ0ZpbHRlcnM9JTVCJTIyeWNkY19wdWJsaWMlMjIlNUQmYW5hbHl0aWNzVGFncz0lNUIlMjJ5Y2RjJTIyJTVE"
	)

	url := fmt.Sprintf("https://%s-dsn.algolia.net/1/indexes/YCCompany_production/query", appID)

	query := map[string]interface{}{
		"hitsPerPage": 1000,
		"query":       "",
		"filters":     fmt.Sprintf("batch:\"%s\"", batchName),
	}

	payload, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Algolia-Application-Id", appID)
	req.Header.Set("X-Algolia-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close response body from %s: %w", url, closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: status code %d for %s", resp.StatusCode, url)
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		err = fmt.Errorf("failed to read response body from %s: %w", url, readErr)
		return nil, err
	}

	var algResp struct {
		Hits []struct {
			Name        string   `json:"name"`
			Slug        string   `json:"slug"`
			Description string   `json:"long_description"`
			Batch       string   `json:"batch"`
			Logo        string   `json:"small_logo_thumb_url"`
			Website     string   `json:"website"`
			Tags        []string `json:"tags"`
			Location    string   `json:"all_locations"`
		} `json:"hits"`
	}

	if err = json.Unmarshal(body, &algResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", url, err)
	}

	for _, h := range algResp.Hits {
		startups = append(startups, Startup{
			Name:        h.Name,
			Slug:        h.Slug,
			Description: h.Description,
			Batch:       h.Batch,
			Logo:        h.Logo,
			Website:     h.Website,
			Tags:        h.Tags,
			Location:    h.Location,
		})
	}

	return startups, err
}
