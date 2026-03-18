# Scry Data Contract

All ingestion scrapers must output data in `.jsonl` (JSON Lines) format. The Go indexer reads this file line-by-line to build the index.

## Schema

Each line must be a valid JSON object with the following fields:

* **`id`** *(string, required)*: A globally unique identifier. Format should be `<source>_<slug>` to prevent database collisions.
* **`source`** *(string, required)*: The name of the dataset or wiki. Used for CLI filtering.
* **`title`** *(string, required)*: The header or title of the entry.
* **`content`** *(string, required)*: The raw text payload to be indexed by SQLite.
* **`url`** *(string, optional)*: The raw hyperlink to the original webpage.
* **`tags`** *(array of strings, optional)*: Additional metadata for targeted searches.

## Example

```json
{
  "id": "sts_wiki_scry",
  "source": "sts_wiki",
  "title": "Scry",
  "content": "Scry is a mechanic unique to the Watcher...",
  "url": "https://slay-the-spire.fandom.com/wiki/Scry",
  "tags": ["mechanic", "watcher"]
}
```
