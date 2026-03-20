import json
import os
import sys

from sources import scrapers


ROOT = "~/.local/share/scry"

def main():
    if len(sys.argv) < 2:
        print("missing argument: expect a source")
        exit(1)

    source = sys.argv[1]

    data_dir_path = os.path.join(os.path.expanduser(ROOT), "data")
    if not os.path.isdir(data_dir_path):
        print("Initializing data path at " + data_dir_path)
        os.makedirs(data_dir_path, exist_ok=False)

    Scraper = scrapers.get(source)
    if Scraper == None:
        print("unknown source: scraper not available")
        exit(1)

    s = Scraper()
    page_contents = s.fetch()
    output_path = os.path.join(data_dir_path, s.name + ".jsonl")
    with open(output_path, "w", encoding="utf-8") as f:
        for page_content in page_contents:
            _ = f.write(json.dumps(page_content) + "\n")


if __name__ == "__main__":
    main()
