package main

import (
	"database/sql"
	"fmt"

	"github.com/manifoldco/promptui"
)

// selectStartup prompts the user to choose a startup from the DB
func selectStartup(db *sql.DB) (*Startup, error) {
	rows, err := db.Query("SELECT name, slug, description, batch, logo, website, tags, location FROM startups ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to query startups: %w", err)
	}
	defer rows.Close()

	var startups []Startup
	for rows.Next() {
		var s Startup
		var tags string
		if err := rows.Scan(&s.Name, &s.Slug, &s.Description, &s.Batch, &s.Logo, &s.Website, &tags, &s.Location); err != nil {
			return nil, fmt.Errorf("failed to scan startup: %w", err)
		}
		startups = append(startups, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating startups: %w", err)
	}
	if len(startups) == 0 {
		return nil, fmt.Errorf("no startups in database")
	}

	names := make([]string, len(startups))
	for i, s := range startups {
		names[i] = s.Name
	}

	prompt := promptui.Select{
		Label: "Select Startup",
		Items: names,
		Size:  20,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &startups[idx], nil
}
