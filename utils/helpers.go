package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SaveResults saves both raw responses and parsed data
func SaveResults(results []*RequestResult, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, result := range results {
		if result.Error != "" || result.Body == "" {
			fmt.Printf("Skipping %s due to error: %s\n", result.ProductID, result.Error)
			continue
		}

		// Determine file extension
		ext := ".html"
		if result.IsJSON {
			ext = ".json"
		}

		// Save raw response
		rawPath := filepath.Join(outputDir, "raw", result.ProductID+ext)
		if err := os.MkdirAll(filepath.Dir(rawPath), 0755); err != nil {
			fmt.Printf("Failed to create directory for %s: %v\n", rawPath, err)
			continue
		}

		if err := os.WriteFile(rawPath, []byte(result.Body), 0644); err != nil {
			fmt.Printf("Failed to save raw file %s: %v\n", rawPath, err)
			continue
		}

		// Save parsed data if available
		if result.ParsedData != nil {
			parsedData, err := json.MarshalIndent(result.ParsedData, "", "  ")
			if err != nil {
				fmt.Printf("Failed to marshal parsed data for %s: %v\n", result.ProductID, err)
				continue
			}

			parsedPath := filepath.Join(outputDir, "parsed", result.ProductID+"_parsed.json")
			if err := os.MkdirAll(filepath.Dir(parsedPath), 0755); err != nil {
				fmt.Printf("Failed to create directory for %s: %v\n", parsedPath, err)
				continue
			}

			if err := os.WriteFile(parsedPath, parsedData, 0644); err != nil {
				fmt.Printf("Failed to save parsed file %s: %v\n", parsedPath, err)
			}
		}

		fmt.Printf("Saved %s (%s)\n", result.ProductID, result.ContentType)
	}

	return nil
}

// Legacy functions (kept for compatibility if needed elsewhere)

// extractProductID extracts product ID from URL (improved version)
func extractProductID(url string) string {
	// This is now handled by the crawler, but keeping for compatibility
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return ""
	}
	last := parts[len(parts)-1]
	if strings.Contains(last, "?") {
		last = strings.SplitN(last, "?", 2)[0]
	}
	if strings.Contains(last, "-") {
		return last
	}
	return ""
}

// saveResponses saves individual responses (deprecated - use SaveResults instead)
func saveResponses(results []*RequestResult, outputDir string) error {
	fmt.Println("Warning: saveResponses is deprecated, use SaveResults instead")
	return SaveResults(results, outputDir)
}

func readURLsFromCSV(filename string) ([]string, error) {
	// Use default if empty
	if filename == "" {
		filename = "../data/urls_to_scrape.csv"
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var urls []string
	for i, row := range records {
		// Skip header
		if i == 0 {
			continue
		}
		if len(row) > 0 && isValidURL(row[0]) { // row[0] is the urls value
			urls = append(urls, row[0])
		}
	}

	return urls, nil
}

func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
