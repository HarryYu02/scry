from typing import override
import re
import requests
from bs4 import BeautifulSoup

from ..base_scraper import BaseScraper


class STSFandomScraper(BaseScraper):
    def __init__(self):
        super().__init__()
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

    def _get_page(self, url: str):
        return requests.get(url, headers=self.headers)

    def _get_all_page_links(self) -> list[str]:
        page = self._get_page(self.site_map)
        output: list[str] = []
        while True:
            soup = BeautifulSoup(page.content, "html.parser")
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
            page = self._get_page(next_url)
        return output

    @override
    def fetch(self):
        links = self._get_all_page_links()
        print(len(links))
        data: list[dict[str, str]] = []
        return data

