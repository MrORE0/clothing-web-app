import asyncio
import json
from CroppScraper import CroppScraper


async def scrape_cropp_website():
    """Main function to scrape the entire Cropp website"""

    scraper = CroppScraper(
        max_concurrent_pages=3,  # less for cloud resources
        context_pool_size=2,
        batch_size=20,  # smaller batches for better memory management
        save_interval=5,
    )

    await scraper.load_progress()

    urls = [
        "https://www.cropp.com/bg/bg/mini-roklya-s-detayl-vazel-458ct-66x",
        "https://www.cropp.com/bg/bg/mini-dress-447ct-99x",
        "https://www.cropp.com/bg/bg/teniska-s-print-127as-05x",
        "https://www.cropp.com/bg/bg/mini-roklya-459ct-mlc",
    ]

    print(f"Starting scrape of {len(urls)} URLs...")

    try:
        results = await scraper.scrape_website_batch(urls)

        print("Scraping completed!")
        print(f"Successfully scraped: {len(results)} products")
        print(f"Failed URLs: {len(scraper.failed_urls)}")

        # Save final results using Pydantic's JSON serialization mode
        with open("cropp_products.json", "w", encoding="utf-8") as f:
            products_json = [product.model_dump(mode="json") for product in results]
            json.dump(products_json, f, indent=2, ensure_ascii=False)

        return results

    except Exception as e:
        print(f"Scraping failed: {str(e)}")
        await scraper.save_progress()  # Save progress even on failure
        raise


if __name__ == "__main__":
    asyncio.run(scrape_cropp_website())
