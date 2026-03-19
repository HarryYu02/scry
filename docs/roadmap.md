# Roadmap

~~1. **XDG compliance:** Update scraper to write to and indexer to read from `~/.local/share/scry/`.~~
2. **Persist indexes:** Split CLI into `index` and `search`. Serialize the TF-IDF index to disk using Go's `encoding/gob` for instant loads.
3. **Multi-wiki expansion:** Write lightweight Python subclasses in the scraper to support new games cleanly.
4. **Smarter search:** Implement basic stemming and Levenshtein distance for fuzzy matching in the Go backend.
5. **TUI integration:** Pipe output into a native pager or use a lightweight library for clean scrolling.
6. **Prepare for release:** Add a Makefile for global binary installation.
