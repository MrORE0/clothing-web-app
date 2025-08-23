package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/MrORE0/clothing-web-app/models"
)

type APIResponse struct {
	Products    []models.RawProduct `json:"products"`
	TotalAmount int                 `json:"productsTotalAmount"`
}

type APIClient struct {
	BaseURL string
	Client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) fetch(ctx context.Context, rawURL string) (*APIResponse, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract brand from domain like "arch.mohito.com" → "mohito.com"
	hostParts := strings.Split(parsedURL.Host, ".")
	refererDomain := ""
	if len(hostParts) >= 2 {
		refererDomain = hostParts[len(hostParts)-2] + "." + hostParts[len(hostParts)-1]
	}

	// Build referer
	referer := fmt.Sprintf("https://www.%s/", refererDomain)

	// Build request
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", referer)

	// Make request
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

	var result APIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &result, nil
}

func (c *APIClient) FetchAllProductsToFile(outputPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testURL := c.BaseURL + "?offset=0&pageSize=1&filters[sortBy]=3&flags[enableddiscountfilter]=true&flags[colorspreviewinfilters]=true&flags[quickshop]=true&flags[loadmorebutton]=1&flags[filterscounter]=1"

	initialResp, err := c.fetch(ctx, testURL)
	if err != nil {
		return fmt.Errorf("failed initial test request: %w", err)
	}
	if initialResp.TotalAmount == 0 {
		return fmt.Errorf("no products found")
	}

	finalURL := fmt.Sprintf("%s?offset=0&pageSize=%d&filters[sortBy]=3&flags[enableddiscountfilter]=true&flags[colorspreviewinfilters]=true&flags[quickshop]=true&flags[loadmorebutton]=1&flags[filterscounter]=1",
		c.BaseURL, initialResp.TotalAmount)

	fullResp, err := c.fetch(ctx, finalURL)
	if err != nil {
		return fmt.Errorf("failed full request: %w", err)
	}

	var products []models.Product

	for _, raw := range fullResp.Products {
		parsed := raw.ToParsed()
		if len(parsed.Variants) > 0 {
			products = append(products, parsed)
		}
	}

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

	fmt.Printf("Successfully scraped %d products to %s\n", len(products), outputPath)
	return nil
}
