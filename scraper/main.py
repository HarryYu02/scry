import asyncio
import json
import os
import sys
import time
from typing import NoReturn

from sources import scrapers, base_scraper


DATA_ROOT = "~/.local/share/scry"
FPS = 30

stats = {
    "count": 0,
    "errors": 0,
    "last_count": 0,
    "start_time": time.time(),
}

def print_error_and_exit(err: str) -> NoReturn:
    print(err, file=sys.stderr)
    sys.exit(1)

async def display_loading_state(s: base_scraper.BaseScraper):
    stats["start_time"] = time.time()
    while True:
        spinning_wheel = "-"
        total_pages = len(s.page_urls)
        completed = stats["count"]
        process_percent = round(completed * 100 / max(total_pages, 1), 1)

        elapsed = time.time() - stats["start_time"]
        req_per_sec = round(completed / elapsed if elapsed > 0 else 0, 1)
        remaining_pages = total_pages - completed
        eta_seconds = remaining_pages / req_per_sec if req_per_sec > 0 else 0
        eta_str = f"{int(eta_seconds // 60):02d}m{int(eta_seconds % 60):02d}s"

        _ = sys.stderr.write(
            f"\r[{spinning_wheel}] Scraping {s.name}... {completed}/{total_pages} ({process_percent}%)" +
            f" | {req_per_sec} req/s" +
            f" | ETA: {eta_str}" +
            f" | Err: {stats['errors']}\033[K"
        )
        _ = sys.stderr.flush()
        stats["last_count"] = completed
        await asyncio.sleep(1 / FPS)

async def async_main():
    if len(sys.argv) < 2:
        print_error_and_exit("ERROR: missing argument: expect a source")

    i = 1
    while sys.argv[i].startswith("-"):
        i += 1
    source = sys.argv[i]

    data_dir_path = os.path.join(os.path.expanduser(DATA_ROOT), "data")
    if not os.path.isdir(data_dir_path):
        print("Initializing data path at " + data_dir_path, file=sys.stderr)
        os.makedirs(data_dir_path, exist_ok=False)

    Scraper = scrapers.get(source)
    if Scraper == None:
        print_error_and_exit("ERROR: unknown source: scraper not available")


    async with Scraper() as s:
        loading_ui = asyncio.create_task(display_loading_state(s))

        def on_success(url: str) -> None:
            stats["count"] += 1

        def on_failed(url: str, err: str) -> None:
            stats["count"] += 1
            stats["errors"] += 1

        page_contents = await s.fetch(on_success=on_success, on_failed=on_failed)
        _ = sys.stderr.write(
            f"\r\033[KScraped {s.name}: {stats['count']} pages successfully" +
            f" ({stats['errors']} failed)\n"
        )

        output_path = os.path.join(data_dir_path, s.name + ".jsonl")
        with open(output_path, "w", encoding="utf-8") as f:
            for page_content in page_contents:
                _ = f.write(json.dumps(page_content) + "\n")

        _ = loading_ui.cancel()


def main():
    asyncio.run(async_main())

if __name__ == "__main__":
    main()
