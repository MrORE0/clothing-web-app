import json
import time
from typing import Set, List
from models.product import Product


class ProgressManager:
    """Manages scraping progress and recovery from failures"""

    def __init__(self, save_interval: int = 10):
        self.save_interval = save_interval
        self.processed_urls: Set[str] = set()
        self.failed_urls: List[str] = []
        self.scraped_products: List[Product] = []

    async def save_progress(self, filename: str = "data/progress.json"):
        """Save current progress to recover from failures"""
        progress_data = {
            "processed_urls": list(self.processed_urls),
            "failed_urls": self.failed_urls,
            "scraped_count": len(self.scraped_products),
            "timestamp": time.time(),
        }

        with open(filename, "w", encoding="utf-8") as f:
            json.dump(progress_data, f, indent=2, ensure_ascii=False)

    async def load_progress(self, filename: str = "data/progress.json"):
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

    def add_processed_url(self, url: str):
        """Mark a URL as processed"""
        self.processed_urls.add(url)

    def add_failed_url(self, url: str):
        """Mark a URL as failed"""
        self.failed_urls.append(url)

    def add_scraped_product(self, product: Product):
        """Add a successfully scraped product"""
        self.scraped_products.append(product)

    def should_save_progress(self) -> bool:
        """Check if progress should be saved based on interval"""
        return len(self.scraped_products) % self.save_interval == 0

    def is_processed(self, url: str) -> bool:
        """Check if a URL has already been processed"""
        return url in self.processed_urls

    def get_stats(self) -> dict:
        """Get current scraping statistics"""
        return {
            "processed": len(self.processed_urls),
            "failed": len(self.failed_urls),
            "scraped": len(self.scraped_products),
        }
