# Roadmap

1. ~~**XDG compliance:** Update scraper to write to and indexer to read from `~/.local/share/scry/`.~~

2. ~~**Persist indexes:** Split CLI into `index` and `search`. Serialize the TF-IDF index to disk using Go's `encoding/gob` for instant loads.~~

3. ~~**Multi-wiki expansion:** Write lightweight Python subclasses in the scraper to support new games cleanly.~~

4. ~~**TUI integration:** Pipe output into a native pager or use a lightweight library for clean scrolling.~~

5. ~~**Prepare for release:** Add a Makefile for global binary installation.~~

6. ~~**Developer UX & Make Targets**: Add a shell wrapper with a shebang and dedicated make targets (like make scrape) so maintainers don't have to manually invoke uv run to test the pipeline.~~

8. ~~**Smarter search:** Implement basic stemming and Levenshtein distance for fuzzy matching in the Go backend.~~

9. ~~**Concurrent Ingestion**: Refactor the Python scraper using asyncio or ThreadPoolExecutor to parallelize document downloads and drastically cut down that 2-hour BG3 scrape time.~~

10. ~~**Scry-Scraper UI Polish:** Implement a minimalist stdout progress spinner/bar for long I/O jobs and clean up terminal logging to prevent raw HTML dumps on failed requests.~~

11. ~~**Stop word list:** Implement stop-word map to drop junk word~~

12. ~~**Output formats:** Implement --stdout, --editor, or other flags for search to format output.~~

13. ~~**Help command:** Implement a more comprehensive help command for each subcommand and their flags.~~

14. ~~**Optimize searching algo:** Optimize the weight of each search term base on document's title.~~

15. ~~**Scry-CLI UI Polish:** Add contextual match snippets with ANSI keyword highlights, strictly respect the `NO_COLOR` environment variable, and structure the `--help` menu output.~~

16. ~~**fzf integration:** Allows live search via fzf scripting.~~

17. ~~**List command:** Implement a list command to view / manage indexes.~~

18. ~~**Shell integration:** Add completion for shell integration.~~

19. ~~**Sync command:** Implement a `sync` command to fetch the latest data.tar.gz from GitHub release page.~~

20. ~~**Rolling Dataset Releases:** Create a `gh` CLI script (or GitHub Action) to automatically upload freshly scraped `.jsonl.gz` datasets to a rolling `latest-datasets` release tag, preventing git history bloat while keeping users up-to-date.~~

21. **Custom Scraper Documentation**: Write a clear tutorial in docs/ detailing how contributors can inherit from the base Python scraper class to add support for their own game wikis.

22. **Structured Logging:** Implement standard library logging (Go `slog`, Python `logging`). Route all logs to `stderr` to keep `stdout` pure for terminal piping.

23. **Comprehensive Testing:** Implement unit and integration tests across all three core components (Python Scraper, Go Indexer, and Go Search CLI) and add a `make test` target to execute them all.

24. **Incremental indexing:** Only index data that has been modified.

25. **Concurrent indexing:** Index multiple documents concurrently to speed up indexing.

26. **Deamon mode:** Run Scry as a daemon / server an a `scryd` binary.

27. **Compress index:** Explore options to compress the index to make the file size smaller. Can sacrifice a little bit of performance.

28. **Scry index loading UI:** Implement a UI for indexing.

29. **Splitting scry-scraper:** Splitting scry-scraper into its own repo.
