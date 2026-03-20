import requests
from bs4 import BeautifulSoup
import html2text

class BaseScraper():
    def __init__(self, name: str, base_url: str, site_map: str, headers: dict[str, str]):
        self.name: str = name
        self.base_url: str = base_url
        self.site_map: str = site_map
        self.headers: dict[str, str] = headers

    def _get_page_soup(self, url: str):
        page = requests.get(url, headers=self.headers)
        soup = BeautifulSoup(page.content, "html.parser")
        return soup

    def _html_to_md(self, html: str) -> str:
        h = html2text.HTML2Text()
        h.body_width = 0
        h.ignore_links = True
        h.ignore_images = True
        return h.handle(html)

    def _get_all_page_links(self) -> list[str]:
        raise NotImplementedError("Scraper must implement _get_all_page_links()")

    def _get_title_from_soup(self, soup: BeautifulSoup) -> str:  # pyright: ignore[reportUnusedParameter]
        raise NotImplementedError("Scraper must implement _get_title_from_soup()")

    def _get_content_from_soup(self, soup: BeautifulSoup) -> str:  # pyright: ignore[reportUnusedParameter]
        raise NotImplementedError("Scraper must implement _get_title_from_soup()")

    def _get_page_content(self, link: str) -> dict[str, str | list[str]] | None:
        soup = self._get_page_soup(link)
        title = self._get_title_from_soup(soup)
        if len(title) == 0:
            return None
        main = self._get_content_from_soup(soup)
        if len(main) == 0:
            return None
        return {
            "id": f'{self.name}_{title.lower().replace(" ", "_")}',
            "source": self.name,
            "title": title,
            "content": self._html_to_md(main),
            "url": link,
            "tags": []
        }

    def _get_all_page_contents(self, links: list[str]) -> list[dict[str, str | list[str]]]:
        output: list[dict[str, str | list[str]]] = []
        for index, link in enumerate(links):
            page_content = self._get_page_content(link)
            if page_content == None:
                print(f"failed #{index}")
                continue
            output.append(page_content)
            print(f"success #{index}")
        return output

    def fetch(self) -> list[dict[str, str | list[str]]]:
        links = self._get_all_page_links()
        print(f"Found {len(links)} pages")
        output = self._get_all_page_contents(links)
        return output
