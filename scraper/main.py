import json

from sources.stsfandom import STSFandomScraper


def main():
    s = STSFandomScraper()
    page_contents = s.fetch()
    with open("stsfandom.jsonl", "w", encoding="utf-8") as f:
        for page_content in page_contents:
            _ = f.write(json.dumps(page_content) + "\n")


if __name__ == "__main__":
    main()
