package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/MrORE0/clothing-web-app/models"
)

var (
	httpClient *http.Client
	once       sync.Once
)

func getHTTPClient(timeout time.Duration) *http.Client {
	once.Do(func() {
		transport := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 30,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  false,
			DisableKeepAlives:   false,
		}

		httpClient = &http.Client{
			Transport: transport,
			Timeout:   timeout,
		}
	})
	return httpClient
}

// UnifiedRequester handles both HTML and JSON requests with the same pattern
type UnifiedRequester struct {
	client *http.Client
	config *models.Config
}

func NewUnifiedRequester(config *models.Config) *UnifiedRequester {
	return &UnifiedRequester{
		client: getHTTPClient(config.Timeout),
		config: config,
	}
}

// ProcessRequests handles both HTML and JSON requests uniformly
func (ur *UnifiedRequester) ProcessRequests(requests []models.RequestInfo) []*models.RequestResult {
	if len(requests) == 0 {
		return nil
	}

	requestChan := make(chan models.RequestInfo, len(requests))
	resultChan := make(chan *models.RequestResult, len(requests))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < ur.config.MaxConcurrency; i++ {
		wg.Add(1)
		go ur.worker(requestChan, resultChan, &wg)
	}

	// Send requests to workers
	go func() {
		for _, req := range requests {
			requestChan <- req
		}
		close(requestChan)
	}()

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	results := make([]*models.RequestResult, 0, len(requests))
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

func (ur *UnifiedRequester) worker(requestChan <-chan models.RequestInfo, resultChan chan<- *models.RequestResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for reqInfo := range requestChan {
		result := ur.makeRequest(reqInfo)
		resultChan <- result
	}
}

func (ur *UnifiedRequester) makeRequest(reqInfo models.RequestInfo) *models.RequestResult {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), ur.config.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", reqInfo.URL, nil)
	if err != nil {
		return &models.RequestResult{
			URL:       reqInfo.URL,
			ProductID: reqInfo.ProductID,
			Duration:  time.Since(start).Milliseconds(),
			Error:     err.Error(),
		}
	}

	// Set headers based on request type
	if reqInfo.IsJSON {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0 Safari/537.36")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Referer", "https://www.cropp.com/bg/bg/")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	} else {
		req.Header.Set("User-Agent", "FastRequester/1.0")
		req.Header.Set("Accept", "*/*")
	}

	// Add custom headers
	for key, value := range reqInfo.Headers {
		req.Header.Set(key, value)
	}

	resp, err := ur.client.Do(req)
	if err != nil {
		return &models.RequestResult{
			URL:       reqInfo.URL,
			ProductID: reqInfo.ProductID,
			Duration:  time.Since(start).Milliseconds(),
			Error:     err.Error(),
		}
	}
	defer resp.Body.Close()

	result := &models.RequestResult{
		URL:         reqInfo.URL,
		ProductID:   reqInfo.ProductID,
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Duration:    time.Since(start).Milliseconds(),
		IsJSON:      reqInfo.IsJSON,
	}

	// Read body if requested
	if ur.config.IncludeBody {
		body, err := io.ReadAll(io.LimitReader(resp.Body, int64(ur.config.MaxBodySize)))
		if err != nil {
			result.Error = fmt.Sprintf("body read error: %v", err)
		} else {
			result.Size = len(body)
			result.Body = string(body)

			// Parse data if requested and it's JSON
			if ur.config.ParseData && reqInfo.IsJSON && len(body) > 0 {
				result.ParsedData = ur.parseJSONResponse(body)
			}
		}
	} else {
		// Just get the size without reading the full body
		written, _ := io.Copy(io.Discard, resp.Body)
		result.Size = int(written)
	}

	return result
}

func (ur *UnifiedRequester) parseJSONResponse(body []byte) interface{} {
	var product models.Product
	if err := json.Unmarshal(body, &product); err != nil {
		// If parsing as Product fails, try to parse as generic interface
		var generic interface{}
		if err := json.Unmarshal(body, &generic); err != nil {
			return fmt.Sprintf("Parse error: %v", err)
		}
		return generic
	}
	return product
}

// BuildJSONURL constructs the Dynamic Yield API URL
// TODO: These parameters likely need to be refreshed periodically
func BuildJSONURL(productID string) string {
	// These parameters should probably be extracted to config or refreshed dynamically
	baseURL := "https://st-eu.dynamicyield.com/spa/json"
	params := url.Values{}
	params.Set("sec", "9879622")
	params.Set("id", "7750695909644877370")
	params.Set("ref", "")
	params.Set("jsession", "dl3sealm7agj3gx6yhg7ppuefdwjnk5a") // TODO: This likely expires
	params.Set("isSesNew", "false")
	params.Set("ctx", fmt.Sprintf(`{"lng":"bg_BG","type":"PRODUCT","data":["%s"]}`, productID))

	return baseURL + "?" + params.Encode()
}
