package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urls = make(map[string]string)

func main() {
    // Initialize random seed
    rand.Seed(time.Now().UnixNano())

    http.HandleFunc("/", handleForm)
    http.HandleFunc("/shorten", handleShorten)
    http.HandleFunc("/short/", handleRedirect)

    fmt.Println("URL Shortener is running on :3030")
    http.ListenAndServe(":3030", nil)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        http.Redirect(w, r, "/shorten", http.StatusSeeOther)
        return
    }

    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>URL Shortener</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    max-width: 800px;
                    margin: 40px auto;
                    padding: 20px;
                    background-color: #f5f5f5;
                }
                .container {
                    background-color: white;
                    padding: 30px;
                    border-radius: 8px;
                    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                }
                h2 {
                    color: #333;
                    text-align: center;
                    margin-bottom: 30px;
                }
                form {
                    display: flex;
                    flex-direction: column;
                    gap: 15px;
                }
                input[type="url"] {
                    padding: 12px;
                    border: 1px solid #ddd;
                    border-radius: 4px;
                    font-size: 16px;
                }
                input[type="submit"] {
                    padding: 12px;
                    background-color: #007bff;
                    color: white;
                    border: none;
                    border-radius: 4px;
                    cursor: pointer;
                    font-size: 16px;
                }
                input[type="submit"]:hover {
                    background-color: #0056b3;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h2>URL Shortener</h2>
                <form method="post" action="/shorten">
                    <input type="url" name="url" placeholder="Enter a URL" required>
                    <input type="submit" value="Shorten">
                </form>
            </div>
        </body>
        </html>
    `)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    originalURL := r.FormValue("url")
    if originalURL == "" {
        http.Error(w, "URL parameter is missing", http.StatusBadRequest)
        return
    }

    shortKey := generateShortKey()
    urls[shortKey] = originalURL
    shortenedURL := fmt.Sprintf("http://localhost:3030/short/%s", shortKey)

    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, `
        <!DOCTYPE html>
        <html>
        <head>
            <title>URL Shortener</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    max-width: 800px;
                    margin: 40px auto;
                    padding: 20px;
                    background-color: #f5f5f5;
                }
                .container {
                    background-color: white;
                    padding: 30px;
                    border-radius: 8px;
                    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                }
                h2 {
                    color: #333;
                    text-align: center;
                    margin-bottom: 30px;
                }
                .url-info {
                    margin: 20px 0;
                    padding: 15px;
                    background-color: #f8f9fa;
                    border-radius: 4px;
                }
                .url-info p {
                    margin: 10px 0;
                    word-break: break-all;
                }
                .url-info a {
                    color: #007bff;
                    text-decoration: none;
                }
                .url-info a:hover {
                    text-decoration: underline;
                }
                .back-button {
                    display: inline-block;
                    margin-top: 20px;
                    padding: 10px 20px;
                    background-color: #6c757d;
                    color: white;
                    text-decoration: none;
                    border-radius: 4px;
                }
                .back-button:hover {
                    background-color: #545b62;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h2>URL Shortened Successfully</h2>
                <div class="url-info">
                    <p><strong>Original URL:</strong> `, originalURL, `</p>
                    <p><strong>Shortened URL:</strong> <a href="`, shortenedURL, `">`, shortenedURL, `</a></p>
                </div>
                <a href="/" class="back-button">Shorten Another URL</a>
            </div>
        </body>
        </html>
    `)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
    shortKey := strings.TrimPrefix(r.URL.Path, "/short/")
    if shortKey == "" {
        http.Error(w, "Shortened key is missing", http.StatusBadRequest)
        return
    }

    originalURL, found := urls[shortKey]
    if !found {
        http.Error(w, "Shortened key not found", http.StatusNotFound)
        return
    }

    http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func generateShortKey() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    const keyLength = 6

    shortKey := make([]byte, keyLength)
    for i := range shortKey {
        shortKey[i] = charset[rand.Intn(len(charset))]
    }
    return string(shortKey)
}
