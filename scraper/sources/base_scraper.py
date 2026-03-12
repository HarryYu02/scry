class BaseScraper():
    def __init__(self):
        return

    def fetch(self) -> list[dict[str, str | list[str]]]:
        raise NotImplementedError("Scraper must implement fetch()")
