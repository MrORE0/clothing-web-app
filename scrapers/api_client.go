package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MrORE0/clothing-web-app/models"
)

type ListOfProducts struct {
	// all the links for all the products
	ProductCodes []models.ProductCode `json:"products"`
	TotalAmount  int                  `json:"productsTotalAmount"`
}

type APIClient struct {
	BaseURL   string
	SecondURL string
	Client    *http.Client
}

func NewAPIClient(baseURL string, secondURL string) *APIClient {
	return &APIClient{
		BaseURL:   baseURL,
		SecondURL: secondURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *APIClient) fetchProducts(ctx context.Context, rawURL string) (*ListOfProducts, string, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract brand from domain like "arch.mohito.com" → "mohito.com"
	hostParts := strings.Split(parsedURL.Host, ".")
	refererDomain := ""
	if len(hostParts) >= 2 {
		refererDomain = hostParts[len(hostParts)-2] + "." + hostParts[len(hostParts)-1]
	}
	var brand string
	if len(hostParts) >= 3 {
		brand = hostParts[len(hostParts)-2]
	}

	// Build referer
	referer := fmt.Sprintf("https://www.%s/", refererDomain)

	// Build request
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", referer)

	// Make request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, brand, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, brand, fmt.Errorf("failed to read response body: %w", err)
	}

	var result ListOfProducts
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, brand, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return &result, brand, nil
}

// The main function that fetches the product info while respecting the status codes so it does not DOS the website
func (c *APIClient) fetchProductInfoWithRetry(ctx context.Context, url string) (models.Product, error) {
	const maxRetries = 2
	wait := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		product, status, headers, err := c.FetchProductInfo(ctx, url)
		if err == nil && status == http.StatusOK {
			return product.ToParsed(), nil
		}

		if status == http.StatusTooManyRequests || status == http.StatusForbidden {
			if ra := headers.Get("Retry-After"); ra != "" {
				if secs, err := strconv.Atoi(ra); err == nil {
					wait = time.Duration(secs) * time.Second
				}
			}
			log.Printf("Got %d on %s, retrying in %v...", status, url, wait)
			select {
			case <-time.After(wait):
				wait *= 2 // exponential backoff
			case <-ctx.Done():
				return models.Product{}, ctx.Err()
			}
			continue
		}

		return models.Product{}, fmt.Errorf("fetch failed (%d): %w", status, err)
	}

	return models.Product{}, fmt.Errorf("max retries exceeded for %s", url)
}

// Will fetched the raw json response from the API and it will return the parsed final version that we need,
// here we fetch the product info
func (c *APIClient) FetchProductInfo(ctx context.Context, rawURL string) (*models.RawProduct, int, http.Header, error) {
	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, 501, nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract brand from domain like "arch.mohito.com" → "mohito.com"
	hostParts := strings.Split(parsedURL.Host, ".")
	refererDomain := ""
	if len(hostParts) >= 2 {
		refererDomain = hostParts[len(hostParts)-2] + "." + hostParts[len(hostParts)-1]
	}
	var brand string
	if len(hostParts) >= 3 {
		brand = hostParts[len(hostParts)-2]
	}

	// Build referer
	referer := fmt.Sprintf("https://www.%s/", refererDomain)

	// Build request
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, 501, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", referer)

	// Make request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, resp.StatusCode, resp.Header, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(rawURL)
		return nil, resp.StatusCode, resp.Header, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, resp.Header, fmt.Errorf("failed to read response body: %w", err)
	}

	var result models.RawProduct
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, resp.StatusCode, resp.Header, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	result.Brand = brand
	return &result, resp.StatusCode, resp.Header, nil
}

// DONE
// Spawns the workers for each website and gives them jobs, then it saves the response in a file
func (c *APIClient) FetchAllProducts() (*[]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	testURL := c.BaseURL + "?offset=0&pageSize=1&filters[sortBy]=3&flags[enableddiscountfilter]=true&flags[colorspreviewinfilters]=true&flags[quickshop]=true&flags[loadmorebutton]=1&flags[filterscounter]=1"

	initialResp, _, err := c.fetchProducts(ctx, testURL)
	if err != nil {
		return nil, fmt.Errorf("failed initial test request: %w", err)
	}
	if initialResp.TotalAmount == 0 {
		return nil, fmt.Errorf("no products found")
	}

	finalURL := fmt.Sprintf("%s?offset=0&pageSize=%d&filters[sortBy]=3&flags[enableddiscountfilter]=true&flags[colorspreviewinfilters]=true&flags[quickshop]=true&flags[loadmorebutton]=1&flags[filterscounter]=1",
		c.BaseURL, initialResp.TotalAmount)

	fullResp, _, err := c.fetchProducts(ctx, finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed full request: %w", err)
	}

	// Channels
	jobs := make(chan string, len(fullResp.ProductCodes)) // product URLs
	results := make(chan models.Product, len(fullResp.ProductCodes))

	var wg sync.WaitGroup
	numWorkers := 5

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for code := range jobs {
				product, err := c.fetchProductInfoWithRetry(ctx, code)
				if err != nil {
					log.Printf("[Worker %d] failed %s: %v", workerID, code, err)
					continue
				}
				results <- product
			}
		}(i)
	}

	for _, code := range fullResp.ProductCodes {
		productURL := c.SecondURL + code.ProductCode
		jobs <- productURL
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var products []models.Product
	for p := range results {
		products = append(products, p)
	}

	fmt.Printf("Successfully scraped %d products ✅\n", len(products))
	return &products, nil
}
