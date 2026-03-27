from typing import override
import re

from bs4 import BeautifulSoup

from ..base_scraper import BaseScraper


class RootFandomScraper(BaseScraper):
    def __init__(self):
        super().__init__(
            "rootfandom",
            "https://root.fandom.com",
            "https://root.fandom.com/wiki/Local_Sitemap",
            {
                "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:148.0) Gecko/20100101 Firefox/148.0",
                "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
                "Accept-Language": "en-US,en;q=0.9",
                "Accept-Encoding": "gzip, deflate, br, zstd",
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
            },
        )

    @override
    async def _get_all_page_urls(self) -> None:
        page_url = self.base_url
        while True:
            soup = await self._get_page_soup(page_url)
            if soup == None:
                return
            content_div = soup.find(id="mw-content-text")
            if content_div == None:
                continue
            urls = content_div.find_all("a")
            for url in urls:
                href = str(url.get("href"))
                if not href.startswith("/wiki/"):
                    continue
                if href.startswith("/wiki/Local_Sitemap"):
                    continue
                async with self.lock:
                    self.page_urls.append(self.base_domain + href)
            next_link = soup.find("a", text=re.compile("^Next page"))
            if next_link == None:
                break
            page_url = self.base_domain + str(next_link.get("href"))

    @override
    def _get_title_from_soup(self, soup: BeautifulSoup) -> str:
        title = soup.find("h1")
        if title == None:
            return ""
        title_text = title.get_text().strip()
        return title_text

    @override
    def _get_content_from_soup(self, soup: BeautifulSoup) -> str:
        main = soup.find(id="content")
        if main == None:
            return ""
        return str(main)
