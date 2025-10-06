package main

import (
	"log"
	"sync"
	"time"

	"github.com/MrORE0/clothing-web-app/scrapers"
	"github.com/MrORE0/clothing-web-app/utils"
)

func main() {
	var wg sync.WaitGroup

	// Initialize headless browser scraper with 5 persistent contexts (one per worker)
	nOfWorkers := 5

	wg.Add(nOfWorkers)

	start := time.Now()

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.cropp.com/api/1099/category/17991/productsWithoutFilters", "https://arch.cropp.com/api/1099/product/")
		products, err := client.FetchAllProducts()
		if err != nil {
			log.Printf("Cropp female error: %v", err)
		}
		if err := utils.SaveProductsToFile("./data/cropp_female_products.json", *products); err != nil {
			log.Printf("Cropp female save error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.cropp.com/api/1099/category/19173/productsWithoutFilters", "https://arch.cropp.com/api/1099/product/")
		products, err := client.FetchAllProducts()
		if err != nil {
			log.Printf("Cropp male error: %v", err)
		}

		if err := utils.SaveProductsToFile("./data/cropp_male_products.json", *products); err != nil {
			log.Printf("Cropp male save error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.housebrand.com/api/1081/category/2879/productsWithoutFilters", "https://arch.housebrand.com/api/1081/product/")
		products, err := client.FetchAllProducts()
		if err != nil {
			log.Printf("Housebrand female error: %v", err)
		}
		if err := utils.SaveProductsToFile("./data/housebrand_female_products.json", *products); err != nil {
			log.Printf("Housebrand female save error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.housebrand.com/api/1081/category/3055/productsWithoutFilters", "https://arch.housebrand.com/api/1081/product/")
		products, err := client.FetchAllProducts()
		if err != nil {
			log.Printf("Housebrand male error: %v", err)
		}
		if err := utils.SaveProductsToFile("./data/housebrand_male_products.json", *products); err != nil {
			log.Printf("Housebrand male save error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.mohito.com/api/1086/category/1983/productsWithoutFilters", "https://arch.mohito.com/api/1086/product/")
		products, err := client.FetchAllProducts()
		if err != nil {
			log.Printf("Mohito error: %v", err)
		}
		if err := utils.SaveProductsToFile("./data/mohito_products.json", *products); err != nil {
			log.Printf("Mohito save error: %v", err)
		}
	}()

	wg.Wait()
	// Kill the chrome and its resources once we are done
	log.Printf("Finished scraping all brands in %vs \n", time.Since(start))
}
