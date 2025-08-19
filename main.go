package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  Single URL: go run main.go <URL>")
		fmt.Println("  Multiple URLs: go run main.go -file <filename>")
		fmt.Println("  Stdin: echo 'url1\\nurl2' | go run main.go -stdin")
		fmt.Println("\nFlags:")
		fmt.Println("  -c <num>     Max concurrent requests (default: CPU cores * 2)")
		fmt.Println("  -timeout <s> Timeout in seconds (default: 10)")
		fmt.Println("  -body        Include response body (default: true)")
		fmt.Println("  -maxbody <n> Max body size in bytes (default: 1MB)")
		fmt.Println("  -output <dir> Output directory (default: ./data)")
		fmt.Println("  -parse       Parse JSON responses into structs (default: true)")
		os.Exit(1)
	}

	config := parseArgs()
	urls, err := getURLs()
	if err != nil {
		fmt.Printf("Error getting URLs: %v\n", err)
		os.Exit(1)
	}

	if len(urls) == 0 {
		fmt.Println("No URLs provided")
		os.Exit(1)
	}

	fmt.Printf("Starting scraper with config: %+v\n", config)

	crawler := NewCrawler(urls[0])
	requester := NewUnifiedRequester(config)

	productURLs, productIDs, err := crawler.GetProductIDs()
	if err != nil {
		fmt.Printf("Error getting product IDs: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d products\n", len(productIDs))

	// Build request info for both HTML and JSON requests
	var allRequests []RequestInfo

	// Add HTML requests
	for i, url := range productURLs {
		allRequests = append(allRequests, RequestInfo{
			URL:       url,
			ProductID: productIDs[i],
			IsJSON:    false,
		})
	}

	// Add JSON requests
	for _, productID := range productIDs {
		allRequests = append(allRequests, RequestInfo{
			URL:       BuildJSONURL(productID),
			ProductID: productID,
			IsJSON:    true,
		})
	}

	fmt.Printf("Processing %d requests (%d HTML + %d JSON)...\n",
		len(allRequests), len(productURLs), len(productIDs))

	// Process all requests
	results := requester.ProcessRequests(allRequests)

	// Save results
	if err := SaveResults(results, config.OutputDir); err != nil {
		fmt.Printf("Error saving results: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Completed! Saved %d results to %s\n", len(results), config.OutputDir)
}

func parseArgs() *Config {
	config := &Config{
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
