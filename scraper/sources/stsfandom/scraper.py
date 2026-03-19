from typing import override
import re
import requests
from bs4 import BeautifulSoup
import html2text

from ..base_scraper import BaseScraper


class STSFandomScraper(BaseScraper):
    def __init__(self):
        super().__init__()
        self.name: str = "stsfandom"
        self.prefix: str = "stsfandom"
        self.base_url: str = "https://slay-the-spire.fandom.com"
        self.site_map: str = "https://slay-the-spire.fandom.com/wiki/Local_Sitemap"
        self.headers: dict[str, str] = {
            "Host": "slay-the-spire.fandom.com",
            "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:148.0) Gecko/20100101 Firefox/148.0",
            "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
            "Accept-Language": "en-US,en;q=0.9",
            "Accept-Encoding": "gzip, deflate, br, zstd",
            "Referer": "https://slay-the-spire.fandom.com/wiki/Slay_the_Spire_Wiki",
            "Sec-GPC": "1",
            "Connection": "keep-alive",
            "Upgrade-Insecure-Requests": "1",
            "Sec-Fetch-Dest": "document",
            "Sec-Fetch-Mode": "navigate",
            "Sec-Fetch-Site": "same-origin",
            "Sec-Fetch-User": "?1",
            "Priority": "u=0, i",
            "Pragma": "no-cache",
            "Cache-Control": "no-cache",
        }

    def _get_page_soup(self, url: str):
        page = requests.get(url, headers=self.headers)
        soup = BeautifulSoup(page.content, "html.parser")
        return soup

    def _get_all_page_links(self) -> list[str]:
        soup = self._get_page_soup(self.site_map)
        output: list[str] = []
        while True:
            links = soup.find_all("a")
            for link in links:
                href = str(link.get("href"))
                if not href.startswith("/wiki/"):
                    continue
                if href.startswith("/wiki/Local_Sitemap"):
                    continue
                output.append(self.base_url + href)
            next_link = soup.find("a", text=re.compile("^Next page"))
            if next_link == None:
                break
            next_url = self.base_url + str(next_link.get("href"))
            soup = self._get_page_soup(next_url)
        return output

    def _html_to_md(self, html: str) -> str:
        h = html2text.HTML2Text()
        h.body_width = 0
        h.ignore_links = True
        h.ignore_images = True
        return h.handle(html)

    def _get_page_content(self, link: str) -> dict[str, str | list[str]]:
        soup = self._get_page_soup(link)
        title = soup.find("h1")
        if title == None:
            return {}
        title_text = title.get_text().strip()
        main = soup.find(id="mw-content-text")
        if main == None:
            return {}
        return {
            "id": f'{self.prefix}_{title_text.lower().replace(" ", "_")}',
            "source": self.prefix,
            "title": title_text,
            "content": self._html_to_md(str(main)),
            "url": link,
            "tags": []
        }

    def _get_all_page_contents(self, links: list[str]) -> list[dict[str, str | list[str]]]:
        output: list[dict[str, str | list[str]]] = []
        for index, link in enumerate(links):
            page_content = self._get_page_content(link)
            if page_content.get("id") == None:
                print(f"failed #{index}")
                continue
            output.append(page_content)
            print(f"success #{index}")
        return output

    @override
    def fetch(self):
        links = self._get_all_page_links()
        print(f"Found {len(links)} pages")
        output = self._get_all_page_contents(links)
        return output

