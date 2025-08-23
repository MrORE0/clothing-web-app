package main

import (
	"log"
	"sync"

	"github.com/MrORE0/clothing-web-app/scrapers"
)

func main() {
	var wg sync.WaitGroup

	// Category IDs
	femaleCategoryURL := "https://arch.cropp.com/api/1099/category/17991/productsWithoutFilters"
	maleCategoryURL := "https://arch.cropp.com/api/1099/category/19173/productsWithoutFilters"

	wg.Add(2)

	// Start female scraping
	go func() {
		defer wg.Done()
		client := scrapers.NewCroppAPIClient(femaleCategoryURL)
		err := client.FetchAllProductsToFile("./data/female_products.json")
		if err != nil {
			log.Printf("Failed to fetch female products: %v", err)
		}
	}()

	// Start male scraping
	go func() {
		defer wg.Done()
		client := scrapers.NewCroppAPIClient(maleCategoryURL)
		err := client.FetchAllProductsToFile("./data/male_products.json")
		if err != nil {
			log.Printf("Failed to fetch male products: %v", err)
		}
	}()

	wg.Wait()
	log.Println("Finished scraping both female and male products.")
}
