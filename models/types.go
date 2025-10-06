package models

import (
	"net/http"
	"time"
)

// --------- RAW STRUCTS (for unmarshaling JSON from Cropp API) ---------

type RawProduct struct {
	ID          int              `json:"id"`
	ProductCode string           `json:"sku"`
	Brand       string           `json:"brand"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Material    string           `json:"material"`
	URL         string           `json:"url"`
	Sizes       []Sizes          `json:"sizes"`
	Prices      Prices           `json:"price"`
	Colors      []RawColorOption `json:"colorOptions"`
	Photos      []Photo          `json:"photos"`
}

type RawColorOption struct {
	ProductCode string       `json:"sku"`
	URL         string       `json:"url"`
	ColorInfo   RawColorInfo `json:"color"`
	IsActive    bool         `json:"isActive"`
	IsBlocked   bool         `json:"isBlocked"`
}

type RawColorInfo struct {
	ColorName string `json:"name"`
	EngName   string `json:"cssName"`
	Photo     string `json:"photo"`
}

type ProductCode struct {
	ProductCode string `json:"sku"`
}

// --------- PARSED STRUCTS (used internally from us ) ---------

type Product struct {
	ID          int           `json:"id"`
	ProductCode string        `json:"sku"`
	Name        string        `json:"name"`
	Prices      Prices        `json:"prices"`
	Photos      []Photo       `json:"images"`
	URL         string        `json:"url"`
	ColorOption []ColorOption `json:"colorOptions"`
	Sizes       []Sizes       `json:"sizes"`
	Brand       string        `json:"brand,omitempty"`
	Description string        `json:"info,omitempty"`
	Material    string        `json:"material"`
	HasDiscount bool          `json:"hasDiscount"`
}

type ColorOption struct {
	ProductCode string `json:"sku"`
	Photo       string `json:"photo"`
	ColorName   string `json:"name"`
	EngName     string `json:"cssName"`
	URL         string `json:"url"`
}

type Prices struct {
	Currency                        string  `json:"currency"`
	RegularPrice                    float64 `json:"regular"`
	FormatedRegularPrice            string  `json:"formattedRegular"`
	FinalPrice                      float64 `json:"finalPrice"`
	FormatedFinalPrice              string  `json:"formattedFinal"`
	AlternativeCurrencyRegularPrice float64 `json:"alternativeCurrencyMobileRegularPrice"`
	FormatedAltRegularPrice         string  `json:"formattedAlternativeCurrencyMobileRegularPrice"`
	AlternativeCurrencyFinalPrice   float64 `json:"alternativeCurrencyMobileFinalPrice"`
	FormatedAltFinalPrice           string  `json:"formattedAlternativeCurrencyMobileFinalPrice"`
	AlternativeCurrency             string  `json:"alternativeCurrency"`
}

type Photo struct {
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Original string `json:"original"`
}

// --------- TRANSFORM FUNCTION ---------

// It will parse the rawProduct from the API response into our json version
func (rp *RawProduct) ToParsed() Product {
	var variants []ColorOption

	for _, color := range rp.Colors {

		// Create variant
		variant := ColorOption{
			ProductCode: color.ProductCode,
			Photo:       color.ColorInfo.Photo,
			ColorName:   color.ColorInfo.ColorName,
			EngName:     color.ColorInfo.EngName,
			URL:         color.URL,
		}

		variants = append(variants, variant)
	}

	return Product{
		ID:          rp.ID,
		ProductCode: rp.ProductCode,
		Name:        rp.Name,
		Description: rp.Description,
		Material:    rp.Material,
		Prices:      rp.Prices,
		Photos:      rp.Photos,
		URL:         rp.URL,
		ColorOption: variants,
		Sizes:       rp.Sizes,
		Brand:       rp.Brand,
	}
}

type Sizes struct {
	Sizes []SizeInfo // not sure if this will need the name of it in the json which is 0,1,2,....
}

type SizeInfo struct {
	SizeName string `json:"sizeName"`
	Stock    bool   `json:"stock"`
}

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
	ParsedData  interface{} `json:"parsed_data,omitempty"`
}

type Config struct {
	MaxConcurrency int
	Timeout        time.Duration
	IncludeBody    bool
	MaxBodySize    int
	OutputDir      string
	ParseData      bool
}

type RequestInfo struct {
	URL       string
	ProductID string
	IsJSON    bool
	Headers   map[string]string
}

type UnifiedRequester struct {
	client *http.Client
	config *Config
}
