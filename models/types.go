package models

import (
	"net/http"
	"strconv"
	"time"
)

// --------- RAW STRUCTS (for unmarshaling JSON from Cropp API) ---------

type RawProduct struct {
	ID          string     `json:"id"`
	ProductCode string     `json:"sku"`
	Name        string     `json:"name"`
	Price       string     `json:"final_price"`
	SalePrice   *string    `json:"price,omitempty"`
	Currency    string     `json:"currency"`
	Images      []string   `json:"img"`
	URL         string     `json:"url"`
	Colors      []RawColor `json:"colorOptions"`
}

type RawColor struct {
	PreviewPhoto string    `json:"previewPhoto"`
	ID           int       `json:"id"`
	SKU          string    `json:"sku"`
	Photo        string    `json:"photo"`
	Images       Images    `json:"images"`
	Name         string    `json:"name"`
	EngName      string    `json:"cssName"`
	URL          string    `json:"url"`
	Prices       RawPrices `json:"prices"`
	Sizes        []Size    `json:"sizes"`
	BackPhoto    string    `json:"backPhoto"`
	IsActive     bool      `json:"isActive"`
	IsBlocked    bool      `json:"isBlocked"`
	IsInStock    bool      `json:"isInStock"`
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
	ID          string   `json:"id"`
	ProductCode string   `json:"sku"`
	Name        string   `json:"name"`
	Price       float64  `json:"final_price"`
	SalePrice   *float64 `json:"price,omitempty"`
	Currency    string   `json:"currency"`
	Images      []string `json:"img"`
	URL         string   `json:"url"`
	Variants    []Color  `json:"colorOptions"`
}

type Color struct {
	Photo    string `json:"photo"`
	Images   Images `json:"images"`
	Name     string `json:"name"`
	EngName  string `json:"cssName"`
	URL      string `json:"url"`
	Prices   Prices `json:"prices"`
	Currency string `json:"currency"`
	InStock  bool   `json:"in_stock"`
	SKU      string `json:"sku"`
	Sizes    []Size `json:"sizes"`
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

type Images struct {
	PreviewPhoto string `json:"previewPhoto"`
	BackPhoto    string `json:"backPhoto"`
}

// --------- TRANSFORM FUNCTION ---------

func (rp *RawProduct) ConvertToProduct() Product {
	var variants []Color

	// Parse fallback base price from raw product
	basePrice, _ := strconv.ParseFloat(rp.Price, 64)

	for _, color := range rp.Colors {
		// Skip inactive or blocked colors
		if !color.IsActive || color.IsBlocked {
			continue
		}

		// Parse color variant prices
		regularPrice, _ := strconv.ParseFloat(color.Prices.Price, 64)
		finalPrice, _ := strconv.ParseFloat(color.Prices.FinalPrice, 64)

		variantPrice := regularPrice
		if variantPrice == 0 {
			variantPrice = basePrice
		}

		prizes := Prices{
			color.Prices.Currency, regularPrice, finalPrice,
			color.Prices.MobileFinalPrice,
			color.Prices.MobileRegularPrice,
			color.Prices.AlternativeCurrencyFinalPrice,
			color.Prices.AlternativeCurrency,
			color.Prices.HasDiscount,
		}

		// Build images list
		images := []string{}
		if color.PreviewPhoto != "" {
			images = append(images, color.PreviewPhoto)
		}
		if color.BackPhoto != "" {
			images = append(images, color.BackPhoto)
		}
		if len(images) == 0 {
			images = rp.Images
		}

		// Filter valid sizes
		var cleanSizes []Size
		for _, s := range color.Sizes {
			if s.Stock && !s.InTransitStock {
				cleanSizes = append(cleanSizes, s)
			}
		}

		// Create variant
		variant := Color{
			Name:     color.Name,
			Images:   Images{color.PreviewPhoto, color.BackPhoto},
			EngName:  color.EngName,
			Photo:    color.Photo, // This is the small color photo
			Sizes:    cleanSizes,
			URL:      color.URL,
			Prices:   prizes,
			Currency: color.Prices.Currency,
			InStock:  color.IsInStock,
			SKU:      color.SKU,
		}

		variants = append(variants, variant)
	}

	return Product{
		ID:          rp.ID,
		Name:        rp.Name,
		ProductCode: rp.ProductCode,
		URL:         rp.URL,
		Images:      rp.Images,
		Variants:    variants,
	}
}

func (rp *RawProduct) ToParsed() Product {
	parsedPrice, _ := strconv.ParseFloat(rp.Price, 64)
	var salePrice *float64
	if rp.SalePrice != nil {
		if p, err := strconv.ParseFloat(*rp.SalePrice, 64); err == nil {
			salePrice = &p
		}
	}

	var variants []Color
	for _, rawColor := range rp.Colors {
		if !rawColor.IsActive || rawColor.IsBlocked {
			continue
		}

		parsedPrices := rawColor.Prices.ToParsed()
		variantPrice := parsedPrices.Price
		if variantPrice == 0 {
			variantPrice = parsedPrice
		}

		var cleanSizes []Size
		for _, s := range rawColor.Sizes {
			if s.Stock && !s.InTransitStock {
				cleanSizes = append(cleanSizes, s)
			}
		}

		images := []string{}
		if rawColor.PreviewPhoto != "" {
			images = append(images, rawColor.PreviewPhoto)
		}
		if rawColor.BackPhoto != "" {
			images = append(images, rawColor.BackPhoto)
		}
		if len(images) == 0 {
			images = rp.Images
		}

		variants = append(variants, Color{
			Name:     rawColor.Name,
			Photo:    rawColor.Photo,
			Images:   rawColor.Images,
			Sizes:    cleanSizes,
			URL:      rawColor.URL,
			Prices:   parsedPrices,
			Currency: parsedPrices.Currency,
			InStock:  rawColor.IsInStock,
			SKU:      rawColor.SKU,
		})
	}

	return Product{
		ID:          rp.ID,
		ProductCode: rp.ProductCode,
		Name:        rp.Name,
		Price:       parsedPrice,
		SalePrice:   salePrice,
		Currency:    rp.Currency,
		Images:      rp.Images,
		URL:         rp.URL,
		Variants:    variants,
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
