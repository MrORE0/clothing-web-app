package main

import (
	"log"
	"sync"

	"github.com/MrORE0/clothing-web-app/scrapers"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(5)

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.cropp.com/api/1099/category/17991/productsWithoutFilters")
		err := client.FetchAllProductsToFile("./data/cropp_female_products.json")
		if err != nil {
			log.Printf("Cropp female error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.cropp.com/api/1099/category/19173/productsWithoutFilters")
		err := client.FetchAllProductsToFile("./data/cropp_male_products.json")
		if err != nil {
			log.Printf("Cropp male error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.housebrand.com/api/1081/category/2879/productsWithoutFilters")
		err := client.FetchAllProductsToFile("./data/housebrand_products.json")
		if err != nil {
			log.Printf("Housebrand error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.mohito.com/api/1086/category/1983/productsWithoutFilters")
		err := client.FetchAllProductsToFile("./data/mohito_female_products.json")
		if err != nil {
			log.Printf("Mohito femlae error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		client := scrapers.NewAPIClient("https://arch.housebrand.com/api/1081/category/3055/productsWithoutFilters")
		err := client.FetchAllProductsToFile("./data/mohito_male_products.json")
		if err != nil {
			log.Printf("Mohito male error: %v", err)
		}
	}()

	wg.Wait()
	log.Println("Finished scraping all brands.")
}
