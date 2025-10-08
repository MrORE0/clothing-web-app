package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/MrORE0/clothing-web-app/models"
	"github.com/MrORE0/clothing-web-app/utils"
)

type ListOfProductCodes struct {
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
    jar, _ := cookiejar.New(nil)
    return &APIClient{
        BaseURL:   baseURL,
        SecondURL: secondURL,
        Client: &http.Client{
            Timeout: 30 * time.Second,
            Jar:     jar,
        },
    }
}

// buildListingURL returns the API listing URL for a given page size.
func (c *APIClient) buildListingURL(pageSize int) string {
	return fmt.Sprintf("%s?offset=0&pageSize=%d&filters[sortBy]=3&flags[enableddiscountfilter]=true&flags[colorspreviewinfilters]=true&flags[quickshop]=true&flags[loadmorebutton]=1&flags[filterscounter]=1",
		c.BaseURL, pageSize)
}

// getTotalProducts queries the listing endpoint with a tiny page size to discover the total amount.
func (c *APIClient) getTotalProducts(ctx context.Context) (int, error) {
	testURL := c.buildListingURL(1)
	initialResp, _, err := c.fetchProductCodes(ctx, testURL)
	if err != nil {
		return 0, fmt.Errorf("failed initial test request: %w", err)
	}
	if initialResp.TotalAmount == 0 {
		return 0, fmt.Errorf("no products found")
	}
	return initialResp.TotalAmount, nil
}

// fetchAllProductCodes returns all product codes by requesting once with the discovered total size.
func (c *APIClient) fetchAllProductCodes(ctx context.Context, total int) ([]models.ProductCode, error) {
	finalURL := c.buildListingURL(total)
	fullResp, _, err := c.fetchProductCodes(ctx, finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed full request: %w", err)
	}
	return fullResp.ProductCodes, nil
}

// spawnWorkers launches request and parse workers and returns channels and waitgroups.
func (c *APIClient) spawnWorkers(ctx context.Context, numWorkers int, rawProducts chan *models.RawProduct, results chan *models.Product) (chan string, *sync.WaitGroup, *sync.WaitGroup) {
	jobs := make(chan string, numWorkers)

	var wgReq sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wgReq.Add(1)
		go func(workerID int) {
			defer wgReq.Done()
			for code := range jobs {
				utils.RandomSleep()
				_, _, err := c.fetchProductInfo(ctx, code, rawProducts)
				if err != nil {
					log.Printf("[Worker %d] failed %s: %v", workerID, code, err)
					continue
				}
			}
		}(i)
	}

	var wgParse sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wgParse.Add(1)
		go func(workerID int) {
			defer wgParse.Done()
			for product := range rawProducts {
				parsedProduct := product.ToParsed()
				results <- &parsedProduct
			}
		}(i)
	}

	return jobs, &wgReq, &wgParse
}

// The function will return a pointer to a slice containing the product codes and the total number
func (c *APIClient) fetchProductCodes(ctx context.Context, rawURL string) (*ListOfProductCodes, string, error) {
	_, _, body, brand, err := utils.MakeJSONRequest(ctx, c.Client, rawURL)
	if err != nil {
		return nil, brand, err
	}

	var result ListOfProductCodes
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, brand, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return &result, brand, nil
}

// Will fetch the raw json response from the API
func (c *APIClient) fetchProductInfo(ctx context.Context, rawURL string, rawProducts chan *models.RawProduct) (int, http.Header, error) {
	statusCode, headers, body, brand, err := utils.MakeJSONRequest(ctx, c.Client, rawURL)
	if err != nil {
		fmt.Println(rawURL)
		return statusCode, headers, err
	}

	var result models.RawProduct
	if err := json.Unmarshal(body, &result); err != nil {
		return statusCode, headers, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	result.Brand = brand
	// pipe the result to be parsed
	rawProducts <- &result
	return statusCode, headers, nil
}

// Spawns the workers for each website and gives them jobs(the item codes to go and request)
// then rawProduct is then parsed and returns a list of the ready to save products
func (c *APIClient) FetchAllParsedProducts() (*[]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	total, err := c.getTotalProducts(ctx)
	if err != nil {
		return nil, err
	}

	productCodes, err := c.fetchAllProductCodes(ctx, total)
	if err != nil {
		return nil, err
	}

	rawProducts := make(chan *models.RawProduct)
	results := make(chan *models.Product, len(productCodes))

	numWorkers := 5
	jobs, wgReq, wgParse := c.spawnWorkers(ctx, numWorkers, rawProducts, results)

	for _, code := range productCodes {
		productURL := c.SecondURL + code.ProductCode
		jobs <- productURL
	}
	close(jobs)

	go func() {
		wgReq.Wait()
		close(rawProducts)
	}()
	go func() {
		wgParse.Wait()
		close(results)
	}()

	var products []models.Product
	for p := range results {
		products = append(products, *p)
	}

	fmt.Printf("Successfully scraped %d products ✅\n.", len(products))
	return &products, nil
}
