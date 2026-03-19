# Scry
Scry is a modular, offline-first, terminal-native search engine written in Go and Python.
This is a mono-repo for both the Python scraper, the Go indexer, and the Go cli.

## Motivation
According to the [Slay the Spire Fandom Wiki](https://slay-the-spire.fandom.com/wiki/Scry):
> Whenever you Scry, you will look at a certain number of cards on top of your draw pile.
> You then may choose to discard any number of them.

You see, I was playing slay the spire on the plane and I need to look up a keyword, but I have no internet.
And sometimes I need to look up the wiki so frequently on some of the more wiki based games like DnDs and strategy games, the web interface literally slows me down.
So I did what any developer would have done: I wrote my own scarper and search engine.

Just like the mechanic, Scry scrapes wiki pages to your filesystem, and let you query against them.
The idea is to be modular that I can write custom scrapers, fine tune the searching algorithm, and use it anywhere.
Now I can take the wiki with me offline, even port it to any front end I want, for any game.

## Architecture
The Scry project is split into three independent components:
1. **Scry Scraper (scraper/):** Downloads wiki pages and compiles them into a `.jsonl` dataset.
2. **Scry Lib (internal/):** An internal library that builds TF-IDF indexes and search from that dataset.
3. **Scry CLI (cmd/scry/):** A CLI front end for user interaction.

The scraped data and indexes are stored at `$HOME/.local/share/scry`.

## Quick Start
There is no release for the app yet, so you will need to clone the repo and run the cli from the root of the project.

### Dependencies
- go >= 1.25
- uv >= 0.9

### 1. Clone the mono-repo
```bash
git clone https://github.com/HarryYu02/scry.git
cd scry
```

### 2. Run the scraper
```bash
cd scraper
uv run main.py
```
It will take a while.

### 3. Run the cli
```bash
cd ..
go run ./cmd/scry search stsfandom "watcher scry"
```

## Usage
Available commands:
- help
- search

## Contributing
### Dependencies
- go >= 1.25
- uv >= 0.9

### 1. Clone the mono-repo
```bash
git clone https://github.com/HarryYu02/scry.git
cd scry
```

### 2. Run the scraper
```bash
cd scraper
mkdir data
uv run main.py
```
It will take a while.

### 3. Run the cli
```bash
cd ..
go run ./cmd/scry search stsfandom "watcher scry"
```

### Submit a pull request
If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

## Resources
- [TF-IDF (Term Frequency-Inverse Document Frequency)](https://en.wikipedia.org/wiki/Tf%E2%80%93idf)
