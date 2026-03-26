import asyncio
from typing import TypedDict, Callable
from urllib.parse import urlparse, urlunparse, quote, unquote
from yarl import URL

import aiohttp
from bs4 import BeautifulSoup
import html2text


class Page(TypedDict):
    id: str
    source: str
    title: str
    content: str
    url: str
    tags: list[str]


class BaseScraper:
    page_urls: list[str] = []
    page_data: dict[str, Page] = {}

    def __init__(
        self, name: str, base_domain: str, base_url: str, headers: dict[str, str]
    ):
        self.name: str = name
        self.base_domain: str = base_domain
        self.base_url: str = base_url
        self.headers: dict[str, str] = headers

        self.lock: asyncio.Lock = asyncio.Lock()
        self.max_concurrency: int = 5
        self.semaphore: asyncio.Semaphore = asyncio.Semaphore(self.max_concurrency)
        self.session: aiohttp.ClientSession

    async def __aenter__(self):
        self.session = aiohttp.ClientSession()
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):  # pyright: ignore[reportUnknownParameterType, reportMissingParameterType]
        await self.session.close()

    def _normalize_url(self, url: str) -> str:
        parsed = urlparse(url)
        normalized_path = quote(unquote(parsed.path))
        parsed = parsed._replace(path=normalized_path)
        return urlunparse(parsed)

    async def _is_page_visited(self, url: str):
        async with self.lock:
            return url in self.page_data.keys()

    async def _fetch_page(self, url: str):
        normalized = URL(self._normalize_url(url), encoded=True)
        async with self.session.get(normalized, headers=self.headers) as res:
            if res.status < 200 or res.status > 299:
                return None
            return await res.text()

    async def _get_page_soup(self, url: str) -> BeautifulSoup | None:
        page = await self._fetch_page(url)
        if page == None:
            return None
        soup = BeautifulSoup(page, "html.parser")
        return soup

    def _html_to_md(self, html: str) -> str:
        h = html2text.HTML2Text()
        h.body_width = 0
        h.ignore_links = True
        h.ignore_images = True
        return h.handle(html)

    async def _get_all_page_urls(self) -> None:
        raise NotImplementedError("Scraper must implement _get_all_page_urls()")

    def _get_title_from_soup(self, soup: BeautifulSoup) -> str:  # pyright: ignore[reportUnusedParameter]
        raise NotImplementedError("Scraper must implement _get_title_from_soup()")

    def _get_content_from_soup(self, soup: BeautifulSoup) -> str:  # pyright: ignore[reportUnusedParameter]
        raise NotImplementedError("Scraper must implement _get_content_from_soup()")

    async def _get_page_content(
        self,
        url: str,
        on_success: Callable[[str], None] | None = None,
        on_failed: Callable[[str, str], None] | None = None,
    ) -> None:
        async with self.semaphore:
            if await self._is_page_visited(url):
                if on_failed:
                    on_failed(url, "already visited")
                return

            soup = await self._get_page_soup(url)
            if soup == None:
                if on_failed:
                    on_failed(url, "page not found")
                return
            title = self._get_title_from_soup(soup)
            if len(title) == 0:
                if on_failed:
                    on_failed(url, "page title not found")
                return
            main = self._get_content_from_soup(soup)
            if len(main) == 0:
                if on_failed:
                    on_failed(url, "page content not found")
                return

            async with self.lock:
                self.page_data[url] = {
                    "id": f"{self.name}_{title.lower().replace(' ', '_')}",
                    "source": self.name,
                    "title": title,
                    "content": self._html_to_md(main),
                    "url": url,
                    "tags": [],
                }
            if on_success:
                on_success(url)

    async def _get_all_page_contents(
        self,
        on_success: Callable[[str], None] | None = None,
        on_failed: Callable[[str, str], None] | None = None,
    ) -> None:
        tasks = [
            asyncio.create_task(self._get_page_content(url, on_success, on_failed))
            for url in self.page_urls
        ]
        if tasks:
            _ = await asyncio.gather(*tasks)

    async def fetch(
        self,
        on_success: Callable[[str], None] | None = None,
        on_failed: Callable[[str, str], None] | None = None,
    ) -> list[Page]:
        await self._get_all_page_urls()
        await self._get_all_page_contents(on_success, on_failed)
        async with self.lock:
            output = list(self.page_data.values())
            return output
