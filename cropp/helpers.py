import asyncio
from selectors_consts import COOKIE_BUTTON_ID, COOKIE_BANNER, COLOR_NAME


async def dismiss_cookie_banner(page):
    """Dismiss cookie banner if present"""
    try:
        await page.locator(COOKIE_BUTTON_ID).click(force=True)
        await page.wait_for_timeout(500)
        await page.locator(COOKIE_BANNER).wait_for(state="hidden", timeout=5000)
        print("[Cookie] Cookie overlay dismissed.")
    except Exception as e:
        print(f"[Cookie] Force click or overlay wait failed: {e}")


async def wait_for_page_load(page, timeout: int = 10000):
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
