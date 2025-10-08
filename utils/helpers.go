package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/MrORE0/clothing-web-app/models"
)

var userAgents = []string{
	// Chrome (Desktop)
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.6613.120 Safari/537.36",
	"Mozilla/5.0 (Windows NT 11.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.83 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.6613.84 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6478.62 Safari/537.36",

	// Chrome (Mobile/Android)
	"Mozilla/5.0 (Linux; Android 14; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.6613.84 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; SM-G996B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.72 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 12; Mi 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6478.127 Mobile Safari/537.36",

	// Chrome on iOS (CriOS)
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/128.0.6613.84 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/127.0.6533.72 Mobile/15E148 Safari/604.1",

	// Safari (Desktop)
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Safari/605.1.15",

	// Safari (iOS/iPadOS)
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",

	// Firefox (Desktop)
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.6; rv:129.0) Gecko/20100101 Firefox/129.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0",

	// Firefox (Mobile)
	"Mozilla/5.0 (Android 14; Mobile; rv:129.0) Gecko/129.0 Firefox/129.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/128.0 Mobile/15E148 Safari/605.1.15",

	// Microsoft Edge (Chromium)
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.6613.84 Safari/537.36 Edg/128.0.2739.42",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.83 Safari/537.36 Edg/127.0.2651.86",

	// Opera (Chromium)
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6478.127 Safari/537.36 OPR/112.0.5197.25",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6478.127 Safari/537.36 OPR/112.0.5197.25",

	// Samsung Internet (Android)
	"Mozilla/5.0 (Linux; Android 14; SAMSUNG SM-S911B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/24.0 Chrome/120.0.6099.144 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 13; SAMSUNG SM-S908B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/23.0 Chrome/116.0.5845.163 Mobile Safari/537.36",

	// UC Browser (Android)
	"Mozilla/5.0 (Linux; U; Android 13; en-US; M2007J3SG) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 UCBrowser/13.4.2.1306 Mobile Safari/537.36",

	// Vivaldi (Chromium)
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6478.62 Safari/537.36 Vivaldi/6.7.3329.31",

	// Brave (Chromium - uses Chrome UA typically, but include a variant)
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.6533.83 Safari/537.36 Brave/127",

	// Bots/Tools
	"curl/8.6.0",
	"Wget/1.21.4 (linux-gnu)",
	"PostmanRuntime/7.40.0",
	"python-requests/2.31.0",
	"Go-http-client/1.1",
}

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

func RandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func RandomSleep() {
	r := rand.Intn(5)
	time.Sleep(time.Duration(r) * time.Second)
}

// DoJSONGet performs a GET request expecting JSON, applying common headers and deriving
// a proper Referer and brand based on the target URL. It returns the HTTP status,
// headers, raw response body, detected brand, and an error if the request failed
// or returned a non-200 status code.
func MakeJSONRequest(ctx context.Context, client *http.Client, rawURL string) (int, http.Header, []byte, string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, "", fmt.Errorf("failed to parse URL: %w", err)
	}

	hostParts := strings.Split(parsedURL.Host, ".")
	refererDomain := ""
	if len(hostParts) >= 2 {
		refererDomain = hostParts[len(hostParts)-2] + "." + hostParts[len(hostParts)-1]
	}
	brand := ""
	if len(hostParts) >= 3 {
		brand = hostParts[len(hostParts)-2]
	}

	referer := fmt.Sprintf("https://www.%s/", refererDomain)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, brand, fmt.Errorf("failed to create request: %w", err)
	}
	// Realistic browser-like headers
	req.Header.Set("User-Agent", RandomUserAgent())
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", referer)
	// Sec-CH hints often seen on Chromium
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")

	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, brand, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return resp.StatusCode, resp.Header, nil, brand, fmt.Errorf("failed to read response body: %w", readErr)
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, resp.Header, body, brand, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	return resp.StatusCode, resp.Header, body, brand, nil
}
