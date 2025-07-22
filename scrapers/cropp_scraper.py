import asyncio
from typing import List
from playwright.async_api import async_playwright, BrowserContext
from models.product import Product
from utils.helpers import dismiss_cookie_banner
from utils.progress_manager import ProgressManager
from utils.data_extractor import DataExtractor
from scrapers.rame_limiter import RateLimiter

from utils.selectors_consts import NAME_SELECTOR


class CroppScraper:
    def __init__(
        self,
        max_concurrent_pages: int = 3,
        context_pool_size: int = 2,
        batch_size: int = 50,
        save_interval: int = 10,  # Save progress every N products
        headless: bool = True,
        user_agent: str = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        request_delay: float = 0.5,
    ):
        self.max_concurrent_pages = max_concurrent_pages
        self.context_pool_size = context_pool_size
        self.batch_size = batch_size
        self.headless = headless
        self.user_agent = user_agent

        self.data_extractor = DataExtractor()

        # Progress tracking
        self.progress_manager = ProgressManager(save_interval=save_interval)

        # Rate limiting
        self.rate_limiter = RateLimiter(request_delay)

        # Browser pool
        self.browser_contexts: List[BrowserContext] = []
        self.semaphore = asyncio.Semaphore(max_concurrent_pages)

    async def initialize_browser_pool(self, playwright):
        """Initialize a pool of browser contexts for reuse"""
        browser = await playwright.chromium.launch(
            headless=self.headless,
            args=[
                "--disable-blink-features=AutomationControlled",
                "--disable-dev-shm-usage",
                "--disable-gpu",
                "--no-sandbox",
                "--disable-setuid-sandbox",
                "--disable-web-security",
                "--disable-features=VizDisplayCompositor",
                "--memory-pressure-off",
                "--max_old_space_size=4096",
            ],
        )

        # Create context pool
        for _ in range(self.context_pool_size):
            context = await browser.new_context(
                user_agent=self.user_agent,
                viewport={"width": 1920, "height": 1080},
                ignore_https_errors=True,
            )
            self.browser_contexts.append(context)

        return browser

    async def get_context(self) -> BrowserContext:
        """Get an available browser context from the pool"""
        context = self.browser_contexts.pop(0)
        self.browser_contexts.append(context)
        return context

    async def get_product_data_(self, product_url: str) -> Product | None:
        async with self.semaphore:
            await self.rate_limiter.wait_if_needed()

            if product_url in self.progress_manager.processed_urls:
                print(f"Skipping already processed: {product_url}")
                return None

            context = await self.get_context()
            page = await context.new_page()

            try:
                # Set shorter timeout for faster failure detection
                await page.goto(
                    product_url, wait_until="domcontentloaded", timeout=15000
                )
                await dismiss_cookie_banner(page)

                # Quick check if product exists
                try:
                    await page.wait_for_selector(NAME_SELECTOR, timeout=5000)
                except Exception as e:
                    print(f"Product not found or page error: {product_url}, {e}")
                    self.progress_manager.failed_urls.append(product_url)
                    return None

                # Extract data
                product_name = await self.data_extractor.extract_product_name(
                    page.locator(NAME_SELECTOR)
                )
                product_variants = await self.data_extractor.extract_product_variants(
                    page
                )

                product_data = Product(name=product_name, variants=product_variants)
                self.progress_manager.add_processed_url(product_url)

                return product_data

            except Exception as e:
                print(f"Error scraping {product_url}: {str(e)}")
                self.progress_manager.add_failed_url(product_url)
                return None

            finally:
                await page.close()

    async def scrape_website_batch(self, urls: List[str]) -> List[Product]:
        """Scrape URLs in batches with progress saving"""
        results = []

        async with async_playwright() as playwright:
            browser = await self.initialize_browser_pool(playwright)

            try:
                # Process URLs in batches
                for i in range(0, len(urls), self.batch_size):
                    batch_urls = urls[i : i + self.batch_size]
                    print(
                        f"Processing batch {i//self.batch_size + 1}/{(len(urls) + self.batch_size - 1)//self.batch_size}"
                    )

                    # Filter out already processed URLs
                    new_urls = [
                        url
                        for url in batch_urls
                        if url not in self.progress_manager.processed_urls
                    ]

                    if not new_urls:
                        print("All URLs in this batch already processed")
                        continue

                    # Process batch
                    tasks = [self.get_product_data_(url) for url in new_urls]
                    batch_results = await asyncio.gather(*tasks, return_exceptions=True)

                    # Filter successful results
                    successful_results = [
                        r for r in batch_results if isinstance(r, Product)
                    ]
                    results.extend(successful_results)
                    for r in successful_results:
                        self.progress_manager.add_scraped_product(r)

                    # Save progress periodically
                    if self.progress_manager.should_save_progress():
                        await self.progress_manager.save_progress()
                        print(
                            f"Progress saved. Scraped: {len(self.progress_manager.scraped_products)}, Failed: {len(self.progress_manager.failed_urls)}"
                        )

                    # Brief pause between batches
                    await asyncio.sleep(1)

            finally:
                await browser.close()
                await self.progress_manager.save_progress()  # Final save

        return results
