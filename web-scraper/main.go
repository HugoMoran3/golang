package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "github.com/gocolly/colly"
		"github.com/gocolly/colly/debug"
)

type Product struct {
    Url, Image, Name, Price string
}

func main() {
    fmt.Println("Starting the scraper...")

    c := colly.NewCollector(
        colly.AllowedDomains("www.scrapingcourse.com"),
        // Enable debugging
        colly.Debugger(&debug.LogDebugger{}),
    )

    var products []Product

    // Add error handling
    c.OnError(func(r *colly.Response, err error) {
        fmt.Printf("Error while scraping: %v\n", err)
        fmt.Printf("Status code: %d\n", r.StatusCode)
    })

    // Log when request is made
    c.OnRequest(func(r *colly.Request) {
        fmt.Printf("Visiting: %s\n", r.URL)
    })

    // Log response received
    c.OnResponse(func(r *colly.Response) {
        fmt.Printf("Received response: %d bytes\n", len(r.Body))
        // Uncomment to see the HTML content
        // fmt.Printf("Response body: %s\n", string(r.Body))
    })

    c.OnHTML("li.product", func(e *colly.HTMLElement) {
        fmt.Println("Found a product element!")

        product := Product{}
        product.Url = e.ChildAttr("a", "href")
        product.Image = e.ChildAttr("img", "src")
        product.Name = e.ChildText(".product-name")
        product.Price = e.ChildText(".price")

        fmt.Printf("Scraped product: %+v\n", product)
        products = append(products, product)
    })

    c.OnScraped(func(r *colly.Response) {
        fmt.Printf("Finished scraping. Found %d products\n", len(products))

        file, err := os.Create("products.csv")
        if err != nil {
            log.Fatalln("Failed to create output CSV file", err)
        }
        defer file.Close()

        writer := csv.NewWriter(file)
        headers := []string{
            "Url",
            "Image",
            "Name",
            "Price",
        }

        if err := writer.Write(headers); err != nil {
            log.Fatalln("Error writing headers:", err)
        }

        for _, product := range products {
            record := []string{
                product.Url,
                product.Image,
                product.Name,
                product.Price,
            }
            if err := writer.Write(record); err != nil {
                log.Fatalln("Error writing record:", err)
            }
        }
        writer.Flush()

        if err := writer.Error(); err != nil {
            log.Fatalln("Error flushing writer:", err)
        }

        fmt.Println("Successfully wrote data to products.csv")
    })

    fmt.Println("Starting to visit the URL...")
    err := c.Visit("https://www.scrapingcourse.com/ecommerce")
    if err != nil {
        log.Fatalf("Error visiting the URL: %v", err)
    }
}
