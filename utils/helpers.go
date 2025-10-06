package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MrORE0/clothing-web-app/models"
)

func SaveProductsToFile(outputPath string, products []models.Product) error {
	if err := os.MkdirAll("./data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal products: %w", err)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write products file: %w", err)
	}

	fmt.Println("Successfully saved products.")
	return nil
}
