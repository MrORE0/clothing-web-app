package models

import (
	"net/http"
	"time"
)

type ProductVariant struct {
	Color           string          `json:"color" bson:"color"`
	Images          []string        `json:"images,omitempty" bson:"images,omitempty"`
	Sizes           map[string]bool `json:"sizes" bson:"sizes"` // e.g. {"XS": false, "S": true}
	URL             string          `json:"url" bson:"url"`     // Assume validated elsewhere
	Price           float64         `json:"price" bson:"price"`
	DiscountedPrice *float64        `json:"discounted_price,omitempty" bson:"discounted_price,omitempty"` // Optional
	Currency        string          `json:"currency" bson:"currency"`
}

type Product struct {
	Name     string           `json:"name" bson:"name"`
	Variants []ProductVariant `json:"variants" bson:"variants"`
}

type Crawler struct {
	BaseURL string
	Client  *http.Client
}

// RequestResult holds information about HTTP request and response
type RequestResult struct {
	URL         string      `json:"url"`
	ProductID   string      `json:"product_id"`
	StatusCode  int         `json:"status_code"`
	ContentType string      `json:"content_type,omitempty"`
	Size        int         `json:"size"`
	Duration    int64       `json:"duration_ms"`
	Error       string      `json:"error,omitempty"`
	IsJSON      bool        `json:"is_json"`
	Body        string      `json:"body,omitempty"`
	ParsedData  interface{} `json:"parsed_data,omitempty"` // Will hold Product struct when parsed
}

// Config holds configuration for the scraper
type Config struct {
	MaxConcurrency int
	Timeout        time.Duration
	IncludeBody    bool
	MaxBodySize    int
	OutputDir      string
	ParseData      bool
}

// RequestInfo holds information about a request to be made
type RequestInfo struct {
	URL       string
	ProductID string
	IsJSON    bool
	Headers   map[string]string
}

// UnifiedRequester handles both HTML and JSON requests with the same pattern
type UnifiedRequester struct {
	client *http.Client
	config *Config
}
