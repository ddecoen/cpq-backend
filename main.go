package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Product represents a product in the catalog
type Product struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	BasePrice   float64 `json:"base_price"`
	PricingType string  `json:"pricing_type"`
	Tiers       []Tier  `json:"tiers,omitempty"`
}

type Tier struct {
	Name        string  `json:"name"`
	MinQuantity int     `json:"min_quantity"`
	MaxQuantity int     `json:"max_quantity"`
	Price       float64 `json:"price"`
}

type PricingResponse struct {
	SKUID           string             `json:"sku_id"`
	ProductName     string             `json:"product_name"`
	Quantity        int                `json:"quantity"`
	TermMonths      int                `json:"term_months"`
	BasePrice       float64            `json:"base_price"`
	Subtotal        float64            `json:"subtotal"`
	Discounts       []DiscountApplied  `json:"discounts"`
	TotalDiscount   float64            `json:"total_discount"`
	FinalPrice      float64            `json:"final_price"`
	AnnualPrice     float64            `json:"annual_price"`
	MonthlyPrice    float64            `json:"monthly_price"`
}

type DiscountApplied struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Percentage  float64 `json:"percentage"`
	Amount      float64 `json:"amount"`
}

type QuoteRequest struct {
	CustomerID string `json:"customer_id"`
	SKUID      string `json:"sku_id"`
	Quantity   int    `json:"quantity"`
	TermMonths int    `json:"term_months,omitempty"`
}

type Quote struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customer_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Items       []QuoteItem `json:"items"`
	Subtotal    float64   `json:"subtotal"`
	TotalDiscount float64 `json:"total_discount"`
	Total       float64   `json:"total"`
}

type QuoteItem struct {
	SKUID       string  `json:"sku_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TermMonths  int     `json:"term_months"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
	Discount    float64 `json:"discount"`
	Total       float64 `json:"total"`
}

type Customer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Company  string `json:"company"`
	Tier     string `json:"tier"`
}

// Enterprise Licenses and AI Add-ons Catalog
var productCatalog = []Product{
	// Enterprise Licenses
	{
		ID:          "sku-1",
		SKU:         "ENT-STARTER",
		Name:        "Enterprise Starter",
		Description: "Perfect for small teams (10 users)",
		Category:    "enterprise_license",
		BasePrice:   50.0,
		PricingType: "concurrent_users",
		Tiers: []Tier{{Name: "10 Users", MinQuantity: 1, MaxQuantity: 10, Price: 50.0}},
	},
	{
		ID:          "sku-2",
		SKU:         "ENT-GROWTH",
		Name:        "Enterprise Growth",
		Description: "Scalable solution (50 users)",
		Category:    "enterprise_license",
		BasePrice:   40.0,
		PricingType: "concurrent_users",
		Tiers: []Tier{{Name: "50 Users", MinQuantity: 11, MaxQuantity: 50, Price: 40.0}},
	},
	{
		ID:          "sku-3",
		SKU:         "ENT-SCALE",
		Name:        "Enterprise Scale",
		Description: "Large organizations (200 users)",
		Category:    "enterprise_license",
		BasePrice:   30.0,
		PricingType: "concurrent_users",
		Tiers: []Tier{{Name: "200 Users", MinQuantity: 51, MaxQuantity: 200, Price: 30.0}},
	},
	{
		ID:          "sku-4",
		SKU:         "ENT-UNLIMITED",
		Name:        "Enterprise Unlimited",
		Description: "Unlimited users with premium support",
		Category:    "enterprise_license",
		BasePrice:   25.0,
		PricingType: "concurrent_users",
		Tiers: []Tier{{Name: "Unlimited", MinQuantity: 201, MaxQuantity: -1, Price: 25.0}},
	},
	// AI Add-ons
	{
		ID:          "sku-ai-1",
		SKU:         "AI-ASSISTANT",
		Name:        "AI Assistant",
		Description: "Intelligent code completion and suggestions",
		Category:    "ai_addon",
		BasePrice:   15.0,
		PricingType: "per_user",
	},
	{
		ID:          "sku-ai-2",
		SKU:         "AI-ANALYTICS",
		Name:        "AI Analytics",
		Description: "Advanced analytics and insights",
		Category:    "ai_addon",
		BasePrice:   25.0,
		PricingType: "per_user",
	},
	{
		ID:          "sku-ai-3",
		SKU:         "AI-SECURITY",
		Name:        "AI Security",
		Description: "AI-powered security scanning",
		Category:    "ai_addon",
		BasePrice:   35.0,
		PricingType: "per_user",
	},
}

var customers = []Customer{
	{ID: "cust-1", Name: "John Doe", Email: "john@acme.com", Company: "Acme Corp", Tier: "enterprise"},
	{ID: "cust-2", Name: "Jane Smith", Email: "jane@startup.io", Company: "Startup Inc", Tier: "startup"},
}

var quotes = make(map[string]*Quote)
var quoteCounter = 1

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	demo := api.PathPrefix("/demo").Subrouter()

	demo.HandleFunc("/products", getProducts).Methods("GET")
	demo.HandleFunc("/pricing", calculatePricing).Methods("GET")
	demo.HandleFunc("/quote", createQuote).Methods("POST")
	demo.HandleFunc("/quotes", listQuotes).Methods("GET")

	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.Use(corsMiddleware)

	fmt.Println("🚀 CPQ Backend starting on port 8080...")
	fmt.Println("📊 Enterprise Licenses: Starter, Growth, Scale, Unlimited")
	fmt.Println("🤖 AI Add-ons: Assistant, Analytics, Security")
	fmt.Println("💰 Discounts: Volume (20%/30%), Multi-year (15%/25%)")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"service": "CPQ Backend API",
		"version": "1.0.0",
		"status":  "running",
		"features": map[string]interface{}{
			"enterprise_licenses": []string{"Starter (10 users)", "Growth (50 users)", "Scale (200 users)", "Unlimited"},
			"ai_addons":           []string{"AI Assistant", "AI Analytics", "AI Security"},
			"discounts": map[string]string{
				"volume":     "20% at $50K, 30% at $100K annual value",
				"multi_year": "15% for 2+ years, 25% for 3+ years",
			},
		},
		"endpoints": map[string]string{
			"products": "/api/v1/demo/products",
			"pricing":  "/api/v1/demo/pricing?sku_id=sku-3&quantity=1&term_months=36",
			"quote":    "/api/v1/demo/quote (POST)",
		},
	}
	json.NewEncoder(w).Encode(response)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products":  productCatalog,
		"customers": customers,
		"success":   true,
	})
}

func calculatePricing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	skuID := r.URL.Query().Get("sku_id")
	quantityStr := r.URL.Query().Get("quantity")
	termMonthsStr := r.URL.Query().Get("term_months")
	customerID := r.URL.Query().Get("customer_id")

	if skuID == "" || quantityStr == "" {
		http.Error(w, "sku_id and quantity are required", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	termMonths := 12
	if termMonthsStr != "" {
		termMonths, err = strconv.Atoi(termMonthsStr)
		if err != nil {
			http.Error(w, "Invalid term_months", http.StatusBadRequest)
			return
		}
	}

	pricing := calculateProductPricing(skuID, quantity, termMonths, customerID)
	if pricing == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(pricing)
}

func calculateProductPricing(skuID string, quantity, termMonths int, customerID string) *PricingResponse {
	product := getProductBySKU(skuID)
	if product == nil {
		return nil
	}

	basePrice := getBasePriceForQuantity(product, quantity)
	subtotal := basePrice * float64(quantity) * float64(termMonths)

	discounts := []DiscountApplied{}
	totalDiscountAmount := 0.0

	// Volume discounts (20% at $50K, 30% at $100K annual value)
	annualValue := basePrice * float64(quantity) * 12
	if annualValue >= 100000 {
		discountAmount := subtotal * 0.30
		discounts = append(discounts, DiscountApplied{
			Type: "volume", Description: "Volume discount: 30% off for $100K+ annual value",
			Percentage: 30.0, Amount: discountAmount,
		})
		totalDiscountAmount += discountAmount
	} else if annualValue >= 50000 {
		discountAmount := subtotal * 0.20
		discounts = append(discounts, DiscountApplied{
			Type: "volume", Description: "Volume discount: 20% off for $50K+ annual value",
			Percentage: 20.0, Amount: discountAmount,
		})
		totalDiscountAmount += discountAmount
	}

	// Multi-year discounts (15% for 2+ years, 25% for 3+ years)
	if termMonths >= 36 {
		discountAmount := subtotal * 0.25
		discounts = append(discounts, DiscountApplied{
			Type: "multi_year", Description: "Multi-year discount: 25% off for 3+ year terms",
			Percentage: 25.0, Amount: discountAmount,
		})
		totalDiscountAmount += discountAmount
	} else if termMonths >= 24 {
		discountAmount := subtotal * 0.15
		discounts = append(discounts, DiscountApplied{
			Type: "multi_year", Description: "Multi-year discount: 15% off for 2+ year terms",
			Percentage: 15.0, Amount: discountAmount,
		})
		totalDiscountAmount += discountAmount
	}

	// Customer tier discounts
	if customerID != "" {
		customer := getCustomerByID(customerID)
		if customer != nil && customer.Tier == "startup" {
			discountAmount := subtotal * 0.10
			discounts = append(discounts, DiscountApplied{
				Type: "customer_tier", Description: "Startup discount: 10% off",
				Percentage: 10.0, Amount: discountAmount,
			})
			totalDiscountAmount += discountAmount
		}
	}

	finalPrice := subtotal - totalDiscountAmount
	annualPrice := finalPrice / float64(termMonths) * 12
	monthlyPrice := finalPrice / float64(termMonths)

	return &PricingResponse{
		SKUID: skuID, ProductName: product.Name, Quantity: quantity, TermMonths: termMonths,
		BasePrice: basePrice, Subtotal: subtotal, Discounts: discounts,
		TotalDiscount: totalDiscountAmount, FinalPrice: math.Round(finalPrice*100)/100,
		AnnualPrice: math.Round(annualPrice*100)/100, MonthlyPrice: math.Round(monthlyPrice*100)/100,
	}
}

func createQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CustomerID == "" || req.SKUID == "" || req.Quantity <= 0 {
		http.Error(w, "customer_id, sku_id, and quantity are required", http.StatusBadRequest)
		return
	}

	if req.TermMonths == 0 {
		req.TermMonths = 12
	}

	customer := getCustomerByID(req.CustomerID)
	if customer == nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	pricing := calculateProductPricing(req.SKUID, req.Quantity, req.TermMonths, req.CustomerID)
	if pricing == nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	quoteID := fmt.Sprintf("quote-%d", quoteCounter)
	quoteCounter++

	quote := &Quote{
		ID: quoteID, CustomerID: req.CustomerID, Status: "draft",
		CreatedAt: time.Now(), ExpiresAt: time.Now().AddDate(0, 0, 30),
		Items: []QuoteItem{{
			SKUID: req.SKUID, ProductName: pricing.ProductName, Quantity: req.Quantity,
			TermMonths: req.TermMonths, UnitPrice: pricing.BasePrice, Subtotal: pricing.Subtotal,
			Discount: pricing.TotalDiscount, Total: pricing.FinalPrice,
		}},
		Subtotal: pricing.Subtotal, TotalDiscount: pricing.TotalDiscount, Total: pricing.FinalPrice,
	}

	quotes[quoteID] = quote

	response := map[string]interface{}{
		"quote": quote, "pricing": pricing, "customer": customer, "success": true,
	}
	json.NewEncoder(w).Encode(response)
}

func listQuotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	customerID := r.URL.Query().Get("customer_id")
	var filteredQuotes []*Quote
	for _, quote := range quotes {
		if customerID == "" || quote.CustomerID == customerID {
			filteredQuotes = append(filteredQuotes, quote)
		}
	}

	response := map[string]interface{}{"quotes": filteredQuotes, "success": true, "count": len(filteredQuotes)}
	json.NewEncoder(w).Encode(response)
}

// Helper functions
func getProductBySKU(skuID string) *Product {
	for _, product := range productCatalog {
		if product.ID == skuID {
			return &product
		}
	}
	return nil
}

func getCustomerByID(customerID string) *Customer {
	for _, customer := range customers {
		if customer.ID == customerID {
			return &customer
		}
	}
	return nil
}

func getBasePriceForQuantity(product *Product, quantity int) float64 {
	if product.PricingType == "concurrent_users" && len(product.Tiers) > 0 {
		for _, tier := range product.Tiers {
			if quantity >= tier.MinQuantity && (tier.MaxQuantity == -1 || quantity <= tier.MaxQuantity) {
				return tier.Price
			}
		}
	}
	return product.BasePrice
}
