import asyncio
import json
from scrapers.scraper import Scraper


async def scrape_website():
    """Main function to scrape the entire website"""

    scraper = Scraper(
        max_concurrent_pages=3,  # less for cloud resources
        context_pool_size=2,
        batch_size=20,  # smaller batches for better memory management
        save_interval=5,
    )

    await scraper.progress_manager.load_progress()

    urls = [
        "https://www.mohito.com/bg/bg/midi-pola-952ee-29p?place=home&brick=new-collection-item-2-3607626",
        "https://www.reserved.com/bg/bg/pulover-sas-sadarzhanie-na-valna-358gp-69x",
        "https://www.cropp.com/bg/bg/kas-suitshart-s-kachulka-921aq-01x",
    ]

    print(f"Starting scrape of {len(urls)} URLs...")

    try:
        results = await scraper.scrape_website_batch(urls)

        print("Scraping completed!")
        print(f"Successfully scraped: {len(results)} products")
        print(f"Failed URLs: {len(scraper.progress_manager.failed_urls)}")

        # Save final results using Pydantic's JSON serialization mode
        with open("data/products.json", "w", encoding="utf-8") as f:
            products_json = [product.model_dump(mode="json") for product in results]
            json.dump(products_json, f, indent=2, ensure_ascii=False)

        return results

    except Exception as e:
        print(f"Scraping failed: {str(e)}")
        await scraper.progress_manager.save_progress()  # Save progress even on failure
        raise


if __name__ == "__main__":
    asyncio.run(scrape_website())
