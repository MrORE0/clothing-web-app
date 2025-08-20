package scrapers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Crawler struct {
	BaseURL string
	Client  *http.Client
}

func NewCrawler(baseURL string) *Crawler {
	return &Crawler{
		BaseURL: baseURL,
		Client:  getHTTPClient(10 * time.Second),
	}
}

// GetProductIDs collects product links and extracts IDs (fixed URL handling)
func (c *Crawler) GetProductIDs() ([]string, []string, error) {
	resp, err := c.Client.Get(c.BaseURL)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	baseURLParsed, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, nil, err
	}

	var urls []string
	var ids []string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && c.isProductLink(attr.Val) {
					id := c.extractID(attr.Val)
					if id != "" {
						// Convert relative URL to absolute
						absoluteURL := c.makeAbsoluteURL(baseURLParsed, attr.Val)
						urls = append(urls, absoluteURL)
						ids = append(ids, id)
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(doc)

	return urls, ids, nil
}

func (c *Crawler) makeAbsoluteURL(base *url.URL, href string) string {
	if strings.HasPrefix(href, "http") {
		return href // Already absolute
	}
	return base.Scheme + "://" + base.Host + href
}

func (c *Crawler) isProductLink(href string) bool {
	// TODO: might need to change the prefix (/bg/bg/) when enrolling in different countries
	return strings.HasPrefix(href, "/bg/bg/") &&
		strings.Contains(href, "-") &&
		!strings.Contains(href, "?")
}

func (c *Crawler) extractID(href string) string {
	parts := strings.Split(href, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}
