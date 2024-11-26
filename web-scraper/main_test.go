// main_test.go
package main

import (
	"os"
	"testing"
)

func TestScrapeProducts(t *testing.T) {
	url := "http://www.scrapingcourse.com/ecommerce/"
	products, err := ScrapeProducts(url)

	if err != nil {
		t.Errorf("Failed to scrape products: %v", err)
	}

	if len(products) == 0 {
		t.Error("No products were scraped")
	}

	// Test first product has required fields
	if len(products) > 0 {
		p := products[0]
		if p.Url == "" {
			t.Error("Product URL is empty")
		}
		if p.Name == "" {
			t.Error("Product name is empty")
		}
		if p.Price == "" {
			t.Error("Product price is empty")
		}
	}
}

func TestSaveToCSV(t *testing.T) {
	// Create test data
	testProducts := []Product{
		{
			Url:   "http://test.com/product1",
			Image: "http://test.com/image1.jpg",
			Name:  "Test Product 1",
			Price: "$10.00",
		},
	}

	// Test filename
	testFile := "test_products.csv"

	// Clean up after test
	defer os.Remove(testFile)

	// Test saving to CSV
	err := SaveToCSV(testProducts, testFile)
	if err != nil {
		t.Errorf("Failed to save to CSV: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("CSV file was not created")
	}
}

func TestInvalidURL(t *testing.T) {
	_, err := ScrapeProducts("http://invalid.url")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}
