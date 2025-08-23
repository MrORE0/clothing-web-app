package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/MrORE0/clothing-web-app/models"
)

type CroppAPIResponse struct {
	Products    []models.RawProduct `json:"products"` // <--- FIXED
	TotalAmount int                 `json:"productsTotalAmount"`
}

type CroppAPIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewCroppAPIClient() *CroppAPIClient {
	return &CroppAPIClient{
		BaseURL: "https://arch.cropp.com/api/1099/category/17991/productsWithoutFilters",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CroppAPIClient) fetch(ctx context.Context, url string) (*CroppAPIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "https://www.cropp.com/")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result CroppAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &result, nil
}

func (c *CroppAPIClient) FetchAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testURL := c.BaseURL + "?offset=0&pageSize=1&filters[sortBy]=3&flags[enableddiscountfilter]=true&flags[colorspreviewinfilters]=true&flags[quickshop]=true&flags[loadmorebutton]=1&flags[filterscounter]=1"

	initialResp, err := c.fetch(ctx, testURL)
	if err != nil {
		return nil, fmt.Errorf("failed initial test request: %w", err)
	}
	if initialResp.TotalAmount == 0 {
		return nil, fmt.Errorf("no products found")
	}

	finalURL := fmt.Sprintf("%s?offset=0&pageSize=%d&filters[sortBy]=3&flags[enableddiscountfilter]=true&flags[colorspreviewinfilters]=true&flags[quickshop]=true&flags[loadmorebutton]=1&flags[filterscounter]=1",
		c.BaseURL, initialResp.TotalAmount)

	fullResp, err := c.fetch(ctx, finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed full request: %w", err)
	}

	var products []models.Product

	for _, raw := range fullResp.Products {
		parsed := raw.ToParsed()
		if len(parsed.Variants) > 0 {
			products = append(products, parsed)
		}
	}

	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	outputPath := "./data/products.json"
	jsonData, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal products: %w", err)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to write products file: %w", err)
	}

	fmt.Printf("Successfully scraped %d products from Cropp\n", len(products))
	return products, nil
}
