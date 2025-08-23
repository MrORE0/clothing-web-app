package models

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// --------- RAW STRUCTS (for unmarshaling JSON from Cropp API) ---------

type RawProduct struct {
	ID          string           `json:"id"`
	ProductCode string           `json:"sku"`
	Name        string           `json:"name"`
	Images      Images           `json:"images"`
	URL         string           `json:"url"`
	Colors      []RawColorOption `json:"colorOptions"`
}

type RawColorOption struct {
	ProductCode string       `json:"sku"`
	URL         string       `json:"url"`
	Sizes       []Size       `json:"sizes"`
	ColorInfo   RawColorInfo `json:"color"`
	IsActive    bool         `json:"isActive"`
	IsBlocked   bool         `json:"isBlocked"`
	Prices      RawPrices    `json:"prices"`
}

type RawColorInfo struct {
	ColorName string `json:"name"`
	EngName   string `json:"cssName"`
	Photo     string `json:"photo"`
}

type RawPrices struct {
	Currency                      string  `json:"currency"`
	Price                         string  `json:"price"`
	FinalPrice                    string  `json:"finalPrice"`
	MobileFinalPrice              float64 `json:"mobileFinalPrice"`
	MobileRegularPrice            float64 `json:"mobileRegularPrice"`
	AlternativeCurrencyFinalPrice string  `json:"alternativeCurrencyFinalPrice"`
	AlternativeCurrency           string  `json:"alternativeCurrency"`
	HasDiscount                   bool    `json:"hasDiscount"`
}

func (rp *RawPrices) ToParsed() Prices {
	price, _ := strconv.ParseFloat(rp.Price, 64)
	finalPrice, _ := strconv.ParseFloat(rp.FinalPrice, 64)
	return Prices{
		Currency:                      rp.Currency,
		Price:                         price,
		FinalPrice:                    finalPrice,
		MobileFinalPrice:              rp.MobileFinalPrice,
		MobileRegularPrice:            rp.MobileRegularPrice,
		AlternativeCurrencyFinalPrice: rp.AlternativeCurrencyFinalPrice,
		AlternativeCurrency:           rp.AlternativeCurrency,
		HasDiscount:                   rp.HasDiscount,
	}
}

// --------- PARSED STRUCTS (used internally from us ) ---------

type Product struct {
	ID          string        `json:"id"`
	ProductCode string        `json:"sku"`
	Name        string        `json:"name"`
	Prices      Prices        `json:"prices"`
	Images      Images        `json:"images"`
	URL         string        `json:"url"`
	Variants    []ColorOption `json:"colorOptions"`
	Sizes       []Size        `json:"sizes"`
}

type ColorOption struct {
	Photo       string `json:"photo"`
	ColorName   string `json:"name"`
	EngName     string `json:"cssName"`
	URL         string `json:"url"`
	ProductCode string `json:"sku"`
}

type Prices struct {
	Currency                      string  `json:"currency"`
	Price                         float64 `json:"price"`
	FinalPrice                    float64 `json:"finalPrice"`
	MobileFinalPrice              float64 `json:"mobileFinalPrice"`
	MobileRegularPrice            float64 `json:"mobileRegularPrice"`
	AlternativeCurrencyFinalPrice string  `json:"alternativeCurrencyFinalPrice"`
	AlternativeCurrency           string  `json:"alternativeCurrency"`
	HasDiscount                   bool    `json:"hasDiscount"`
}

type ImageResolution struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type Images struct {
	Small ImageResolution `json:"850"`
	Large ImageResolution `json:"1200"`
}

// --------- TRANSFORM FUNCTION ---------

func (rp *RawProduct) ToParsed() Product {
	var variants []ColorOption
	var prices Prices
	var sizes []Size

	for _, color := range rp.Colors {

		// Create variant
		variant := ColorOption{
			ColorName:   color.ColorInfo.ColorName,
			EngName:     color.ColorInfo.EngName,
			Photo:       color.ColorInfo.Photo, // This is the small color photo
			URL:         color.URL,
			ProductCode: color.ProductCode,
		}

		// Parse fallback base price from raw product
		basePrice, _ := strconv.ParseFloat(color.Prices.Price, 64)

		// Parse color variant prices
		multiplier := math.Pow(10, float64(2))
		color.Prices.Price = strings.ReplaceAll(color.Prices.Price, ",", ".") // need this in order for the conversion bellow to work
		color.Prices.FinalPrice = strings.ReplaceAll(color.Prices.FinalPrice, ",", ".")

		regularPrice, _ := strconv.ParseFloat(color.Prices.Price, 32)
		regularPrice = math.Round(regularPrice*multiplier) / multiplier

		finalPrice, _ := strconv.ParseFloat(color.Prices.FinalPrice, 32)
		finalPrice = math.Round(finalPrice*multiplier) / multiplier

		variantPrice := regularPrice
		if variantPrice == 0 {
			variantPrice = basePrice
		}

		prices = Prices{
			color.Prices.Currency, regularPrice, finalPrice,
			color.Prices.MobileFinalPrice,
			color.Prices.MobileRegularPrice,
			color.Prices.AlternativeCurrencyFinalPrice,
			color.Prices.AlternativeCurrency,
			color.Prices.HasDiscount,
		}
		variants = append(variants, variant)
		sizes = color.Sizes
	}

	return Product{
		ID:          rp.ID,
		ProductCode: rp.ProductCode,
		Name:        rp.Name,
		Prices:      prices,
		Images:      rp.Images,
		URL:         rp.URL,
		Variants:    variants,
		Sizes:       sizes,
	}
}

type Size struct {
	SizeName       string `json:"sizeName"`
	Stock          bool   `json:"stock"`
	MagentoID      int    `json:"magentoId"`
	SizeID         int    `json:"sizeId"`
	SKU            string `json:"sku"`
	Key            string `json:"key"`
	InTransitStock bool   `json:"inTransitStock"`
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
