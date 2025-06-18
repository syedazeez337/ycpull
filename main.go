package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter" // For formatted table output

	_ "github.com/mattn/go-sqlite3" // SQLite driver (already in database.go but good for explicitness if main is run alone)
)

const dbPath = "yc_startups.db"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <batch_name>")
		fmt.Println("Example: go run . summer-2023")
		os.Exit(1)
	}
	batchName := os.Args[1]

	log.Printf("Fetching data for batch: %s", batchName)
	startups, err := FetchBatchData(batchName)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	if len(startups) == 0 {
		log.Printf("No startups found for batch %s. The API might have returned an empty list or the batch name is incorrect.", batchName)
		// os.Exit(0) // Decide if exiting here is desired or if we should proceed to init DB anyway
	}
	log.Printf("Fetched %d startups from API.", len(startups))

	log.Printf("Initializing database: %s", dbPath)
	db, err := InitDB(dbPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()
	log.Println("Database initialized successfully.")

	log.Println("Storing startups in the database...")
	err = StoreStartups(db, startups)
	if err != nil {
		log.Fatalf("Error storing startups: %v", err)
	}
	log.Println("Startups stored successfully.")

	log.Println("Displaying startups from database:")
	err = DisplayStartups(db)
	if err != nil {
		log.Fatalf("Error displaying startups: %v", err)
	}
}

// DisplayStartups queries and prints the startup data in a table format.
func DisplayStartups(db *sql.DB) error {
	rows, err := db.Query("SELECT name, website, location FROM startups ORDER BY name")
	if err != nil {
		return fmt.Errorf("failed to query startups: %w", err)
	}
	defer rows.Close()

	fmt.Println("\n--- YC Startups ---")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight) // Removed tabwriter.Debug
	fmt.Fprintln(w, "Name\tWebsite\tLocation\t")
	fmt.Fprintln(w, "----\t-------\t--------\t")

	var count int
	for rows.Next() {
		var name, website, location string
		if err := rows.Scan(&name, &website, &location); err != nil {
			return fmt.Errorf("failed to scan startup row: %w", err)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", name, website, location)
		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error during rows iteration: %w", err)
	}

	w.Flush()

	if count == 0 {
		fmt.Println("No startups found in the database to display.")
	} else {
		fmt.Printf("\nDisplayed %d startups.\n", count)
	}

	return nil
}
