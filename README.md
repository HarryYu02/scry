# Scry
Scry is a modular, offline-first, terminal-native search engine written in Go and Python.
This is a mono-repo for both the Python scraper, the Go indexer, and the Go cli.

![Search result of using scry search "shadowheart"](https://github.com/HarryYu02/scry/blob/main/assets/search.png "Search result")
![Rendering search result in less](https://github.com/HarryYu02/scry/blob/main/assets/render.png "Rendering in less")

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

### 1. Installation
Simply download the binary from the latest [release](https://github.com/HarryYu02/scry/releases).
Optionally, you can clone the main branch for the latest version as well.

you will need Go installed on your system:
```bash
git clone https://github.com/HarryYu02/scry.git
cd scry
make install
```

It will install to `$HOME/.local/bin`, make sure `$HOME/.local/bin` is in your $PATH.

### 2. Getting wiki data
You can always write your custom scraper class and scrape the data,
but I have included some example datasets in the [release](https://github.com/HarryYu02/scry/releases/tag/v2.0.0).

Download the attached data.tar.gz files from the release.
Then unzip and move the downloaded datasets to your local data directory:
```bash
tar xzf data.tar.gz
mkdir -p ~/.local/share/scry/data
mv *.jsonl ~/.local/share/scry/data/
```

Optionally, you can install the Python scraper and scrape the data yourself.
You will need uv installed for `scry-scraper`:
```bash
make install-scraper
scry-scraper stsfandom
```

### 3. Try it out
```bash
scry help
scry index stsfandom
scry search stsfandom "watcher scry"
```

## Usage
Available commands:
- help
- index
- search

Use "scry help <command>" for more information about a command.

## Contributing
Please see [contribute](https://github.com/HarryYu02/scry/blob/main/docs/contribute.md) for how to contribute.

## Resources
- [TF-IDF (Term Frequency-Inverse Document Frequency)](https://en.wikipedia.org/wiki/Tf%E2%80%93idf)
- [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance)
- [Damerau–Levenshtein distance](https://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance)
- [Stemming](https://en.wikipedia.org/wiki/Stemming)
- [Snowball](https://github.com/snowballstem/snowball/tree/master)
- [ANSI escape code](https://en.wikipedia.org/wiki/ANSI_escape_code)
