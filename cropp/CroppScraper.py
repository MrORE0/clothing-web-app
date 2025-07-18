import asyncio
import json
import time
from typing import List, Dict, Set
from playwright.async_api import async_playwright, BrowserContext
from product import Product, ProductVariant
from helpers import dismiss_cookie_banner

# Selectors
SIZE_PICKER_SELECTOR = "ul[data-testid*='product-size-group'] li"
NAME_SELECTOR = "h1[data-testid*='product-name']"
COLOR_NAME = "span[data-testid*='color-picker-color-name']"
COLORS = "li[data-testid='color-picker-list-item'] button[data-testid='color-picker-list-item-button']"
PRICE = "div[data-selen='product-price']"
DISCOUNT = "div[data-selen='product-discount-price']"
REGULAR_PRICE = "div[data-selen='product-regular-price']"
COOKIE_BUTTON_ID = "#cookiebotDialogOkButton"


class CroppScraper:
    def __init__(
        self,
        max_concurrent_pages: int = 3,
        context_pool_size: int = 2,
        batch_size: int = 50,
        save_interval: int = 10,  # Save progress every N products
        headless: bool = True,
        user_agent: str = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    ):
        self.max_concurrent_pages = max_concurrent_pages
        self.context_pool_size = context_pool_size
        self.batch_size = batch_size
        self.save_interval = save_interval
        self.headless = headless
        self.user_agent = user_agent

        # Progress tracking
        self.processed_urls: Set[str] = set()
        self.failed_urls: List[str] = []
        self.scraped_products: List[Product] = []

        # Rate limiting
        self.request_delay = 0.5  # Delay between requests
        self.last_request_time = 0

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

    async def respect_rate_limit(self):
        """Implement rate limiting to avoid being blocked/restricted"""
        current_time = time.time()
        time_since_last = current_time - self.last_request_time

        if time_since_last < self.request_delay:
            await asyncio.sleep(self.request_delay - time_since_last)

        self.last_request_time = time.time()

    async def save_progress(self, filename: str = "cropp_progress.json"):
        """Save current progress to recover from failures"""
        progress_data = {
            "processed_urls": list(self.processed_urls),
            "failed_urls": self.failed_urls,
            "scraped_count": len(self.scraped_products),
            "timestamp": time.time(),
        }

        with open(filename, "w", encoding="utf-8") as f:
            json.dump(progress_data, f, indent=2, ensure_ascii=False)

    async def load_progress(self, filename: str = "cropp_progress.json"):
        """Load previous progress"""
        try:
            with open(filename, "r", encoding="utf-8") as f:
                progress_data = json.load(f)
                self.processed_urls = set(progress_data.get("processed_urls", []))
                self.failed_urls = progress_data.get("failed_urls", [])
                print(
                    f"Loaded progress: {len(self.processed_urls)} URLs already processed"
                )
        except FileNotFoundError:
            print("No previous progress found, starting fresh")

    async def get_product_data_(self, product_url: str) -> Product | None:
        """version of get_product_data"""
        async with self.semaphore:
            await self.respect_rate_limit()

            if product_url in self.processed_urls:
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
                except Exception:
                    print(f"Product not found or page error: {product_url}")
                    self.failed_urls.append(product_url)
                    return None

                # Extract data
                product_name = await self.extract_product_name(
                    page.locator(NAME_SELECTOR)
                )
                product_variants = await self.extract_product_variants_(page)

                product_data = Product(name=product_name, variants=product_variants)
                self.processed_urls.add(product_url)

                return product_data

            except Exception as e:
                print(f"Error scraping {product_url}: {str(e)}")
                self.failed_urls.append(product_url)
                return None

            finally:
                await page.close()

    async def extract_product_variants_(self, page) -> List[ProductVariant]:
        """Improved variant extraction with better duplicate prevention"""
        variants = []
        processed_colors = set()  # Track processed color names only
        processed_urls = set()  # Track processed URLs for additional safety

        try:
            # Wait for initial elements
            await asyncio.wait_for(
                asyncio.gather(
                    page.wait_for_selector(COLORS, timeout=8000),
                    page.wait_for_selector(COLOR_NAME, timeout=8000),
                ),
                timeout=10,
            )

            color_count = await page.locator(COLORS).count()
            max_colors = min(color_count, 10)

            print(f"Found {color_count} colors, processing {max_colors}")

            for i in range(max_colors):
                try:
                    # Re-query color buttons to avoid stale references
                    color_buttons = page.locator(COLORS)
                    color_button = color_buttons.nth(i)

                    # Get current color name before clicking
                    try:
                        current_color_name = (
                            await page.locator(COLOR_NAME).nth(0).inner_text()
                        )
                    except Exception:
                        current_color_name = "unknown"

                    print(
                        f"Clicking color button {i+1}/{max_colors}, current color: {current_color_name}"
                    )

                    # Click the color button
                    await color_button.click()

                    # Wait for changes to take effect
                    await self.wait_for_color_change_(page, current_color_name)

                    # Extract new color name and URL
                    try:
                        new_color_name = (
                            await page.locator(COLOR_NAME).nth(0).inner_text()
                        )
                        current_url = page.url
                    except Exception as e:
                        print(f"Failed to extract color name or URL: {e}")
                        continue

                    print(f"Color changed to: {new_color_name}, URL: {current_url}")

                    # Check for duplicates using both color name and URL
                    if new_color_name in processed_colors:
                        print(f"Duplicate color name detected: {new_color_name}")
                        continue

                    if current_url in processed_urls:
                        print(f"Duplicate URL detected: {current_url}")
                        continue

                    # Mark as processed
                    processed_colors.add(new_color_name)
                    processed_urls.add(current_url)

                    # Extract variant data
                    item_sizes = await self.extract_size_availability_(
                        page.locator(SIZE_PICKER_SELECTOR)
                    )
                    price_data = await self.extract_product_price_(page)

                    if price_data:
                        price, discounted_price, currency = price_data
                        variants.append(
                            ProductVariant(
                                color=new_color_name,
                                sizes=item_sizes,
                                url=current_url,
                                price=price,
                                discounted_price=discounted_price,
                                currency=currency,
                            )
                        )
                        print(f"Successfully processed color: {new_color_name}")
                    else:
                        print(
                            f"Failed to extract price data for color: {new_color_name}"
                        )

                except Exception as e:
                    print(f"Error processing color variant {i}: {str(e)}")
                    continue

        except Exception as e:
            print(f"Error extracting variants: {str(e)}")

        print(f"Total variants extracted: {len(variants)}")
        return variants

    async def wait_for_color_change_(
        self, page, previous_color_name: str, timeout: int = 15000
    ):
        """Wait for color change to complete with better verification"""
        try:
            start_time = time.time()

            # Wait for any network activity to settle first
            await asyncio.sleep(0.3)

            # Try to wait for load state with shorter timeout
            try:
                await asyncio.wait_for(
                    page.wait_for_load_state("domcontentloaded"), timeout=3
                )
            except asyncio.TimeoutError:
                pass  # Continue anyway

            # Wait for color name element to be stable and check for change
            max_attempts = 30
            stable_count = 0
            last_color_name = previous_color_name

            for _ in range(max_attempts):
                try:
                    # Check if we've exceeded total timeout
                    if (time.time() - start_time) * 1000 > timeout:
                        break

                    current_color_name = (
                        await page.locator(COLOR_NAME).nth(0).inner_text()
                    )

                    if current_color_name != previous_color_name:
                        # Color has changed, wait for it to stabilize
                        if current_color_name == last_color_name:
                            stable_count += 1
                            if stable_count >= 3:  # Color has been stable for 3 checks
                                print(
                                    f"Color change detected: '{previous_color_name}' → '{current_color_name}'"
                                )
                                await asyncio.sleep(0.2)  # Small buffer
                                return
                        else:
                            stable_count = 0
                            last_color_name = current_color_name

                    await asyncio.sleep(0.1)

                except Exception:
                    await asyncio.sleep(0.1)
                    continue

            print(
                f"Warning: Color change not detected after {max_attempts} attempts (timeout: {timeout}ms)"
            )
            # Even if we didn't detect change, wait a bit for page to stabilize
            await asyncio.sleep(0.5)

        except Exception as e:
            print(f"[Color Change] Error waiting for color change: {e}")
            await asyncio.sleep(0.5)

    async def wait_for_new_page_load_(self, page, timeout: int = 10000):
        """Enhanced page load waiting with better verification"""
        try:
            # Wait for network to be mostly idle
            await asyncio.wait_for(
                page.wait_for_load_state("networkidle", timeout=timeout),
                timeout=timeout / 1000,
            )

            # Ensure essential elements are present and stable
            await asyncio.wait_for(
                page.wait_for_selector(COLOR_NAME, timeout=timeout),
                timeout=timeout / 1000,
            )

            # Additional wait for dynamic content to load
            await asyncio.sleep(0.3)

        except asyncio.TimeoutError:
            print("[Page] Timeout waiting for page load, continuing...")
            # Fallback to shorter wait
            await asyncio.sleep(1)

    async def extract_product_price_(
        self, page
    ) -> tuple[float, float | None, str] | None:
        """price extraction with faster timeouts"""
        try:

            def parse_price(text: str) -> tuple[float, str]:
                temp = text.replace("\xa0", " ").strip().split(" ")
                return float(temp[0].replace(",", ".")), temp[1]

            # Try discount price first (shorter timeout)
            try:
                await asyncio.wait_for(
                    asyncio.gather(
                        page.wait_for_selector(DISCOUNT, timeout=3000),
                        page.wait_for_selector(REGULAR_PRICE, timeout=3000),
                    ),
                    timeout=4,
                )

                discount_text = await page.locator(DISCOUNT).nth(0).inner_text()
                regular_text = await page.locator(REGULAR_PRICE).nth(0).inner_text()

                price, currency = parse_price(regular_text)
                discounted_price, _ = parse_price(discount_text)

                return price, discounted_price, currency

            except Exception:
                # Fallback to regular price
                await page.wait_for_selector(PRICE, timeout=3000)
                raw_price_text = await page.locator(PRICE).nth(0).inner_text()
                price, currency = parse_price(raw_price_text)
                return price, None, currency

        except Exception as e:
            print(f"Price extraction failed: {str(e)}")
            return None

    async def extract_size_availability_(self, picker_items) -> Dict[str, bool]:
        """size extraction"""
        item_sizes = {}
        try:
            count = await picker_items.count()
            # Limit size processing to avoid excessive loops
            max_sizes = min(count, 15)

            for i in range(max_sizes):
                try:
                    item = picker_items.nth(i)
                    label = await item.inner_text()
                    item_class = await item.get_attribute("class")

                    if item_class and label:
                        item_sizes[label] = "inactive" not in item_class.split()

                except Exception:
                    continue

        except Exception as e:
            print(f"Size extraction error: {str(e)}")

        return item_sizes

    async def extract_product_name(self, name_locator) -> str:
        """Extract product name with error handling"""
        try:
            if await name_locator.count() != 1:
                print("Warning: Multiple or no product name elements found")
                return "Unknown Product"

            name_tag = name_locator.nth(0)
            return await name_tag.inner_text()

        except Exception as e:
            print(f"Name extraction error: {str(e)}")
            return "Unknown Product"

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
                        url for url in batch_urls if url not in self.processed_urls
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
                    self.scraped_products.extend(successful_results)

                    # Save progress periodically
                    if len(self.scraped_products) % self.save_interval == 0:
                        await self.save_progress()
                        print(
                            f"Progress saved. Scraped: {len(self.scraped_products)}, Failed: {len(self.failed_urls)}"
                        )

                    # Brief pause between batches
                    await asyncio.sleep(1)

            finally:
                await browser.close()
                await self.save_progress()  # Final save

        return results
