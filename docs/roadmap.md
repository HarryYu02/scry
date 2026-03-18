## Roadmap

1. **System Installation & XDG Compliance:** Update Go to read from `~/.local/share/scry/`. Add a Makefile for global binary installation.
2. **The Package Manager Pivot:** Split CLI into `build` and `search`. Serialize the TF-IDF index to disk using Go's `encoding/gob` for instant loads.
3. **Multi-Wiki Expansion:** Write lightweight Python subclasses in the scraper to support new games cleanly.
4. **Smarter Search Math:** Implement basic stemming and Levenshtein distance for fuzzy matching in the Go backend.
5. **TUI Integration:** Pipe output into a native pager or use a lightweight library for clean scrolling.
