import asyncio
import time
from typing import List, Dict
from models.product import ProductVariant
from utils.selectors_consts import (
    REGULAR_PRICE,
    PRICE,
    DISCOUNT,
    COLORS,
    SIZE_PICKER_SELECTOR,
    get_color_name_selector,
)


class DataExtractor:
    """Handles extraction of product data from web pages"""

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

    async def extract_product_price(
        self, page
    ) -> tuple[float, float | None, str] | None:
        """Enhanced price extraction with faster timeouts"""
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

    async def extract_size_availability(self, picker_items) -> Dict[str, bool]:
        """Extract size availability information"""
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

    async def extract_product_variants(self, page) -> List[ProductVariant]:
        """Extract all product variants with better duplicate prevention"""
        variants = []
        processed_colors = set()  # Track processed color names only
        processed_urls = set()  # Track processed URLs for additional safety

        try:
            # Wait for initial elements
            await asyncio.wait_for(
                asyncio.gather(
                    page.wait_for_selector(COLORS, timeout=8000),
                    page.wait_for_selector(
                        get_color_name_selector(page.url), timeout=8000
                    ),
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
                            await page.locator(get_color_name_selector(page.url))
                            .nth(0)
                            .inner_text()
                        )
                    except Exception:
                        current_color_name = "unknown"

                    print(
                        f"Clicking color button {i+1}/{max_colors}, current color: {current_color_name}"
                    )

                    # Click the color button
                    await color_button.click()

                    # Wait for changes to take effect
                    await self.wait_for_color_change(page, current_color_name)

                    # Extract new color name and URL
                    try:
                        new_color_name = (
                            await page.locator(get_color_name_selector(page.url))
                            .nth(0)
                            .inner_text()
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
                    item_sizes = await self.extract_size_availability(
                        page.locator(SIZE_PICKER_SELECTOR)
                    )
                    price_data = await self.extract_product_price(page)

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

    async def wait_for_color_change(
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
                        await page.locator(get_color_name_selector(page.url))
                        .nth(0)
                        .inner_text()
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
