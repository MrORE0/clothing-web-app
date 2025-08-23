package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/MrORE0/clothing-web-app/models"
	"github.com/MrORE0/clothing-web-app/scrapers"
)

func main() {
	client := scrapers.NewCroppAPIClient()
	_, err := client.FetchAllProducts() // This returns the products but I am currently not using it
	if err != nil {
		log.Fatalf("Failed to fetch products: %v", err)
	}
}

// TODO: This will not be needed with the new approach
func parseArgs() *models.Config {
	config := &models.Config{
		MaxConcurrency: runtime.NumCPU() * 2,
		Timeout:        10 * time.Second,
		IncludeBody:    true,
		MaxBodySize:    1024 * 1024, // 1MB
		OutputDir:      "./data",
		ParseData:      true,
	}

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-c":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &config.MaxConcurrency)
				i++
			}
		case "-timeout":
			if i+1 < len(args) {
				var seconds int
				fmt.Sscanf(args[i+1], "%d", &seconds)
				config.Timeout = time.Duration(seconds) * time.Second
				i++
			}
		case "-body":
			config.IncludeBody = true
		case "-maxbody":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &config.MaxBodySize)
				i++
			}
		case "-output":
			if i+1 < len(args) {
				config.OutputDir = args[i+1]
				i++
			}
		case "-parse":
			config.ParseData = true
		}
	}

	return config
}
