# Roadmap

~~1. **XDG compliance:** Update scraper to write to and indexer to read from `~/.local/share/scry/`.~~

~~2. **Persist indexes:** Split CLI into `index` and `search`. Serialize the TF-IDF index to disk using Go's `encoding/gob` for instant loads.~~

~~3. **Multi-wiki expansion:** Write lightweight Python subclasses in the scraper to support new games cleanly.~~

~~4. **TUI integration:** Pipe output into a native pager or use a lightweight library for clean scrolling.~~

~~5. **Prepare for release:** Add a Makefile for global binary installation.~~

6. **Developer UX & Make Targets**: Add a shell wrapper with a shebang and dedicated make targets (like make scrape) so maintainers don't have to manually invoke uv run to test the pipeline.

7. **Smarter search:** Implement basic stemming and Levenshtein distance for fuzzy matching in the Go backend.

8. **Concurrent Ingestion**: Refactor the Python scraper using asyncio or ThreadPoolExecutor to parallelize document downloads and drastically cut down that 2-hour BG3 scrape time.

9. **Custom Scraper Documentation**: Write a clear tutorial in docs/ detailing how contributors can inherit from the base Python scraper class to add support for their own game wikis.

10. **Rolling Dataset Releases:** Create a `gh` CLI script (or GitHub Action) to automatically upload freshly scraped `.jsonl.gz` datasets to a rolling `latest-datasets` release tag, preventing git history bloat while keeping users up-to-date.
