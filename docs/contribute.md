# Contributing to Scry

First off, thanks for taking the time to contribute!

Scry is built strictly around the Unix philosophy: do one thing, do it well, and pass raw data through clean interfaces.
We prioritize speed, minimalism, and a native terminal experience.

## The Repo Split (Where do I put my code?)
Scry is divided into two separate projects to keep the core engine lightweight. Please ensure you are opening your Pull Request in the right place:

* **[scry](https://github.com/HarryYu02/scry):** You are here. This repository is pure Go. It contains the `scry` CLI client, the `scryd` background daemon, the `bbolt` database logic, and the TF-IDF math.
* **[scry-scraper](https://github.com/HarryYu02/scry-scraper):** This repository is pure Python. It handles all the messy HTML parsing, web scraping, and data distribution. If a wiki is missing data, or you want to add a completely new game to the ecosystem, open a PR there.

## Repository Architecture
This project follows a standard, modular Go layout:
* `cmd/scry/`: The lightweight CLI client. It parses user flags, queries the daemon, and formats the output for the terminal.
* `cmd/scryd/`: The background daemon. It holds the exclusive file lock on the database and listens for queries via a Unix Domain Socket.
* `internal/`: The core engine. All shared logic (lexing, stemming, storage, and rendering) lives here.

## Local Development Setup
1. Clone the repository: `git clone https://github.com/HarryYu02/scry.git`
2. Ensure you have a recent version of Go installed.
3. Run `make` (or `go build -o bin/scry cmd/scry/main.go`) to compile.
4. To populate your local test environment with data, use the built-in sync command to pull the latest `.tar.gz` release: `bin/scry sync`

## Pull Request Guidelines
We prefer quiet, asynchronous collaboration. You don't need to open an issue before submitting a PR for small fixes, but please keep the following in mind:

* **Avoid External Dependencies:** We rely on the Go standard library and `bbolt`. Do not introduce new third-party packages unless absolutely unavoidable.
* **Keep it Fast:** Scry is designed for instant, local retrieval. PRs that introduce noticeable latency or heavy memory footprints will be rejected.
* **Format Your Code:** Always run `go fmt ./...` before committing.
* **Scope Your Changes:** Keep PRs focused on a single feature or bug fix.
