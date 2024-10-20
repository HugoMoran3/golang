package main

import (
    "encoding/csv"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    // Define command-line flags
    inputFile := flag.String("input", "", "Input CSV file")
    outputFile := flag.String("output", "", "Output JSON file")
    flag.Parse()

    // Validate input and output files
    if *inputFile == "" || *outputFile == "" {
        fmt.Println("Please provide both input and output file names.")
        flag.PrintDefaults()
        os.Exit(1)
    }

    // Read CSV file
    csvFile, err := os.Open(*inputFile)
    if err != nil {
        fmt.Println("Error opening CSV file:", err)
        os.Exit(1)
    }
    defer csvFile.Close()

    reader := csv.NewReader(csvFile)
    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error reading CSV:", err)
        os.Exit(1)
    }

    // Convert CSV to JSON
    var result []map[string]string
    headers := records[0]

    for _, record := range records[1:] {
        row := make(map[string]string)
        for i, value := range record {
            row[headers[i]] = value
        }
        result = append(result, row)
    }

    // Write JSON to file
    jsonData, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        fmt.Println("Error marshaling JSON:", err)
        os.Exit(1)
    }

    err = os.WriteFile(*outputFile, jsonData, 0644)
    if err != nil {
        fmt.Println("Error writing JSON file:", err)
        os.Exit(1)
    }

    fmt.Printf("Successfully converted %s to %s\n", *inputFile, *outputFile)
}
