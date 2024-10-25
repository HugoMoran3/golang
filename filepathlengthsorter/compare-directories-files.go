package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: program <dir1> <dir2>")
        os.Exit(1)
    }

    dir1 := os.Args[1]
    dir2 := os.Args[2]

    // Create maps to store filenames
    files1 := make(map[string]bool)
    files2 := make(map[string]bool)

    // Walk first directory
    err := filepath.Walk(dir1, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            relPath, err := filepath.Rel(dir1, path)
            if err != nil {
                return err
            }
            files1[relPath] = true
        }
        return nil
    })
    if err != nil {
        fmt.Printf("Error walking directory %s: %v\n", dir1, err)
        os.Exit(1)
    }

    // Walk second directory
    err = filepath.Walk(dir2, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            relPath, err := filepath.Rel(dir2, path)
            if err != nil {
                return err
            }
            files2[relPath] = true
        }
        return nil
    })
    if err != nil {
        fmt.Printf("Error walking directory %s: %v\n", dir2, err)
        os.Exit(1)
    }

    // Find unique files in first directory
    fmt.Printf("\nFiles only in %s:\n", dir1)
    uniqueCount1 := 0
    for file := range files1 {
        if !files2[file] {
            fmt.Println(file)
            uniqueCount1++
        }
    }

    // Find unique files in second directory
    fmt.Printf("\nFiles only in %s:\n", dir2)
    uniqueCount2 := 0
    for file := range files2 {
        if !files1[file] {
            fmt.Println(file)
            uniqueCount2++
        }
    }

    // Print summary
    fmt.Printf("\nSummary:\n")
    fmt.Printf("Files only in %s: %d\n", dir1, uniqueCount1)
    fmt.Printf("Files only in %s: %d\n", dir2, uniqueCount2)
    fmt.Printf("Total unique files: %d\n", uniqueCount1+uniqueCount2)
}
