# Mobile EPUB Reader

A Go module designed for mobile apps that provides EPUB parsing and HTML content
extraction capabilities. Compatible with `gomobile` for creating Android (AAR)
and iOS (Framework) bindings.

## Features

- ✅ Parse EPUB files from byte arrays
- ✅ Extract book metadata (title, author, description, etc.)
- ✅ Convert chapters to mobile-friendly HTML
- ✅ Navigation functions (next/prev chapter)
- ✅ Table of contents extraction
- ✅ Reading progress tracking
- ✅ Dark mode CSS support
- ✅ gomobile compatible API

## Installation

```bash
go get github.com/Ppang0405/bookshelf/mobile-epub-reader
```

## Usage in Go

```go
package main

import (
    "fmt"
    "os"
    
    epubreader "github.com/Ppang0405/bookshelf/mobile-epub-reader"
)

func main() {
    // Load EPUB file
    epubData, err := os.ReadFile("book.epub")
    if err != nil {
        panic(err)
    }
    
    // Open the EPUB
    err = epubreader.OpenEpubFromBytes(epubData)
    if err != nil {
        panic(err)
    }
    
    // Get book information
    bookInfo := epubreader.GetBookInfo()
    fmt.Println("Book Info:", bookInfo)
    
    // Get first chapter HTML
    chapterHTML := epubreader.GetChapterHTMLByIndex(0)
    fmt.Printf("Chapter HTML length: %d\n", len(chapterHTML))
    
    // Navigation
    if epubreader.HasNextChapterByIndex(0) {
        nextHTML := epubreader.GetNextChapterHTMLByIndex(0)
        fmt.Printf("Next chapter HTML length: %d\n", len(nextHTML))
    }
    
    // Clean up
    epubreader.CloseBook()
}
```

## Building for Mobile

### Prerequisites

```bash
go install golang.org/x/mobile/cmd/gomobile@latest
gomobile init
```

### Build Android AAR

```bash
gomobile bind -target=android -o epubreader.aar github.com/Ppang0405/bookshelf/mobile-epub-reader
```

### Build iOS Framework

```bash
gomobile bind -target=ios -o EpubReader.xcframework github.com/Ppang0405/bookshelf/mobile-epub-reader
```

## Mobile Integration

### Android (Kotlin/Java)

```kotlin
// Load EPUB file
val epubBytes = assets.open("book.epub").readBytes()

// Open EPUB
Epubreader.openEpubFromBytes(epubBytes)

// Get book info (JSON string)
val bookInfo = Epubreader.getBookInfo()
val bookData = JSONObject(bookInfo)
val title = bookData.getString("title")

// Get chapter HTML
val chapterHTML = Epubreader.getChapterHTMLByIndex(0)

// Display in WebView
webView.loadData(chapterHTML, "text/html", "UTF-8")

// Navigation
val hasNext = Epubreader.hasNextChapterByIndex(0)
if (hasNext) {
    val nextHTML = Epubreader.getNextChapterHTMLByIndex(0)
    webView.loadData(nextHTML, "text/html", "UTF-8")
}

// Clean up
Epubreader.closeBook()
```

### iOS (Swift)

```swift
import EpubReader

// Load EPUB file
guard let epubData = NSDataAsset(name: "book")?.data else { return }

// Open EPUB
do {
    try EpubreaderOpenEpubFromBytes(epubData)
} catch {
    print("Error opening EPUB: \(error)")
    return
}

// Get book info
let bookInfo = EpubreaderGetBookInfo()
if let bookData = bookInfo.data(using: .utf8),
   let json = try? JSONSerialization.jsonObject(with: bookData) as? [String: Any] {
    let title = json["title"] as? String
}

// Get chapter HTML
let chapterHTML = EpubreaderGetChapterHTMLByIndex(0)

// Display in WKWebView
webView.loadHTMLString(chapterHTML, baseURL: nil)

// Navigation
let hasNext = EpubreaderHasNextChapterByIndex(0)
if hasNext {
    let nextHTML = EpubreaderGetNextChapterHTMLByIndex(0)
    webView.loadHTMLString(nextHTML, baseURL: nil)
}

// Clean up
EpubreaderCloseBook()
```

## API Reference

### Core Functions

- `OpenEpubFromBytes(data []byte) error` - Load EPUB from byte array
- `CloseBook()` - Close current book and free memory

### Metadata Functions

- `GetBookInfo() string` - Get book metadata as JSON
- `GetTableOfContents() string` - Get TOC as JSON
- `GetChapterCount() int` - Get total number of chapters

### Content Functions

- `GetChapterHTMLByIndex(index int) string` - Get chapter HTML
- `GetChapterTitleByIndex(index int) string` - Get chapter title

### Navigation Functions

- `GetNextChapterHTMLByIndex(currentIndex int) string` - Get next chapter HTML
- `GetPrevChapterHTMLByIndex(currentIndex int) string` - Get previous chapter
  HTML
- `HasNextChapterByIndex(currentIndex int) bool` - Check if next chapter exists
- `HasPrevChapterByIndex(currentIndex int) bool` - Check if previous chapter
  exists

### Progress Functions

- `GetReadingProgress(currentChapter int) float64` - Get reading progress
  (0-100%)

## CSS Styling

The generated HTML includes mobile-friendly CSS with:

- Responsive typography
- Dark mode support (`@media (prefers-color-scheme: dark)`)
- Optimized line height and spacing
- Image scaling
- Clean blockquote styling

## Dependencies

- `github.com/taylorskalyo/goreader/epub` - EPUB parsing
- `golang.org/x/net/html` - HTML processing

## License

Same as the parent bookshelf project.
