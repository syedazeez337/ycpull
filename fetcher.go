package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// No os import needed for this approach
)

// FetchBatchData fetches startup data for a given batch name from the YC API.
func FetchBatchData(batchName string) (startups []Startup, err error) { // Named return error 'err'
	url := fmt.Sprintf("https://ycombinator-oss.vercel.app/api/batch/%s.json", batchName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %w", url, err)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil && err == nil { // If no other error has occurred, assign closeErr to err
			err = fmt.Errorf("failed to close response body from %s: %w", url, closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		// Important: We must return here. If we don't, 'err' will be nil,
		// and the defer func might overwrite it with a potential resp.Body.Close() error,
		// masking the more critical StatusCode error.
		return nil, fmt.Errorf("failed to fetch data: status code %d for %s", resp.StatusCode, url)
	}

	body, readErr := ioutil.ReadAll(resp.Body) // Use new variable for error to avoid conflict with named return 'err'
	if readErr != nil {
		// Assign to 'err' directly if it's the primary error we want to return
		err = fmt.Errorf("failed to read response body from %s: %w", url, readErr)
		return nil, err
	}

	// Unmarshal into the named return variable 'startups'
	unmarshalErr := json.Unmarshal(body, &startups)
	if unmarshalErr != nil {
		// Assign to 'err' directly
		err = fmt.Errorf("failed to unmarshal JSON from %s: %w", url, unmarshalErr)
		return nil, err
	}

	// If err is still nil here, the defer func might set it if resp.Body.Close() fails.
	// Otherwise, any error encountered above (status, read, unmarshal) will be returned.
	return startups, err
}
