# Scry Scraper

This directory contains the web scrapers used to generate the raw text datasets (`.jsonl`) that the Scry Go engine indexes.

## Prerequisites

This project uses [uv](https://github.com/astral-sh/uv) for fast, deterministic Python environment management.

Ensure `uv` and Python 3.13+ are installed on your system.

## Usage

To run a scraper, pass the namespace of the wiki you want to download:

```bash
uv run main.py bg3wiki
uv run main.py stsfandom
```
uv will automatically resolve dependencies from the uv.lock file,
scrape the target site, and output a raw JSONL file to the system data directory:
* Output Path: ~/.local/share/scry/data/<namespace>.jsonl

## Adding a New Wiki

1. Create a new directory under sources/ (e.g., sources/mywiki/).
2. Implement a scraper.py class that inherits from BaseScraper.
3. Ensure you exclude .mw-redirect classes in your CSS selectors to avoid indexing typo/redirect pages.
4. Register the new source in sources/\_\_init\_\_.py.
