import asyncio
import time


class RateLimiter:
    """Handles rate limiting for web scraping to avoid being blocked"""

    def __init__(self, request_delay: float = 0.5):
        self.request_delay = request_delay
        self.last_request_time = 0

    async def wait_if_needed(self):
        """Implement rate limiting by waiting if requests are too frequent"""
        current_time = time.time()
        time_since_last = current_time - self.last_request_time

        if time_since_last < self.request_delay:
            await asyncio.sleep(self.request_delay - time_since_last)

        self.last_request_time = time.time()

    def set_delay(self, delay: float):
        """Update the delay between requests"""
        self.request_delay = delay
