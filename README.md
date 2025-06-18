# YC Startup Fetcher

This Go program fetches startup data batch-wise from Y Combinator's public Algolia index, stores it in a local SQLite database, and prints a table of the startups.

## Features

- Fetches startup data via Algolia using the `YCCompany_production` index.
- Stores data locally in an SQLite database (`yc_startups.db`).
- Displays startup name, website, and location in a formatted table.
- Handles duplicate entries by ignoring conflicts based on startup slug.

## Prerequisites

- Go (version 1.18 or later recommended)

## Setup & Running

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone <repository_url>
    cd <repository_directory>
    ```

2.  **Tidy dependencies:**
    The necessary dependencies (`github.com/mattn/go-sqlite3`) are listed in the `go.mod` file. They will be downloaded automatically when you build or run the program. You can also fetch them explicitly:
    ```bash
    go mod tidy
    ```
    or
    ```bash
    go get .
    ```

3.  **Build the program (optional):**
    You can build an executable:
    ```bash
    go build -o yc_fetcher
    ```
    This will create an executable file named `yc_fetcher`.

4.  **Run the program:**
    You need to provide a batch name as a command-line argument. For example, to fetch data for the 'summer-2023' batch:

    If you built the executable:
    ```bash
    ./yc_fetcher summer-2023
    ```

    Alternatively, you can run directly using `go run`:
    ```bash
    go run . summer-2023
    ```

    Other example batch names: `winter-2023`, `winter-2022`, `summer-2022`, etc.
    The batch names must match the format used on YC's site, e.g. `Summer 2023`.

## Database

- The program will automatically create an SQLite database file named `yc_startups.db` in the same directory where you run the command.
- If you run the program multiple times with different batch names, the new data will be added to the same database. If a startup from a new batch pull already exists (based on its unique 'slug'), it will not be re-added.

## Project Structure

- `main.go`: Contains the main application logic, including command-line argument parsing and orchestration of fetching and storage. Also includes the display logic.
- `models.go`: Defines the `Startup` data structure.
- `fetcher.go`: Contains the logic for fetching data from the YC API.
- `database.go`: Contains the logic for database initialization and data storage.
- `go.mod`, `go.sum`: Go module files.
- `yc_startups.db`: (Created after first run) The SQLite database file.
