class BaseScraper():
    def __init__(self):
        return

    def fetch(self) -> list[dict[str, str]]:
        raise NotImplementedError("Scraper must implement fetch()")
