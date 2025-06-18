package main

// Startup represents the data structure for a single startup.
type Startup struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Batch       string   `json:"batch"`
	Logo        string   `json:"logo"`
	Website     string   `json:"website"`
	Tags        []string `json:"tags"`
	Location    string   `json:"location"`
}
