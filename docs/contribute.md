# Contribute
## Dependencies
- go >= 1.25
- uv >= 0.9

## Submit a pull request
If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

### 1. Fork and clone the mono-repo
```bash
gh repo fork https://github.com/HarryYu02/scry.git
git clone https://github.com/<username>/scry.git
cd scry
```

### 2. Run the scraper
```bash
cd scraper
uv run main.py stsfandom
```
It will take a while (about 5 minutes on my laptop).
See [scry-scraper](https://github.com/HarryYu02/scry/blob/main/scraper/README.md) for how to implement your custom scraper.

### 3. Run the cli
```bash
cd ..
go run ./cmd/scry index stsfandom
go run ./cmd/scry search stsfandom "watcher scry"
```
