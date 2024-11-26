package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

type Product struct {
	Url, Image, Name, Price string
}

func main() {
    // create new collector object instance
    c := colly.NewCollector(
	colly.AllowedDomains("www.scrapingcourse.com"),
	)

    // initialise the slices that will contain scraped data
    var products []Product

    // OnHTML callback
    c.OnHTML(".product", func(e *colly.HTMLElement){

    })

    // called before an HTTP request is triggered
    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting: ", r.URL)
    })

    // triggered when the scraper encounters an error
    c.OnError(func(_ *colly.Response, err error) {
        fmt.Println("Something went wrong: ", err)
    })

    // fired when the server responds
    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Page visited: ", r.Request.URL)
    })

    // triggered when a CSS selector matches an element
    c.OnHTML("a", func(e *colly.HTMLElement) {
        // printing all URLs associated with the <a> tag on the page
        fmt.Println("%v", e.Attr("href"))
    })

    // triggered once scraping is done (e.g., write the data to a CSV file)
    c.OnScraped(func(r *colly.Response) {
        fmt.Println(r.Request.URL, " scraped!")
    })


    // open the target URL
    c.Visit("https://www.scrapingcourse.com/ecommerce")


}
