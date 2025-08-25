package epubreader

import (
	"fmt"
	"os"
	"testing"
)

// Simple test to verify the EPUB reader works
func TestEpubReader(t *testing.T) {
	// This test would need an actual EPUB file
	// For now, we'll just test the basic structure
	
	// Test with empty data (should fail gracefully)
	err := OpenEpubFromBytes([]byte{})
	if err == nil {
		t.Error("Expected error with empty data, got nil")
	}
	
	// Test global reader functions when no book is loaded
	bookInfo := GetBookInfo()
	if bookInfo == "" {
		t.Error("GetBookInfo should return error message when no book loaded")
	}
	
	chapterCount := GetChapterCount()
	if chapterCount != 0 {
		t.Error("Expected 0 chapters when no book loaded")
	}
	
	fmt.Println("Basic tests passed")
}

// Example usage function for documentation
func ExampleUsage() {
	// Load EPUB from file
	epubData, err := os.ReadFile("example.epub")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	
	// Open the EPUB
	err = OpenEpubFromBytes(epubData)
	if err != nil {
		fmt.Printf("Error opening EPUB: %v\n", err)
		return
	}
	
	// Get book information
	bookInfo := GetBookInfo()
	fmt.Printf("Book Info: %s\n", bookInfo)
	
	// Get table of contents
	toc := GetTableOfContents()
	fmt.Printf("TOC: %s\n", toc)
	
	// Get first chapter HTML
	firstChapterHTML := GetChapterHTMLByIndex(0)
	fmt.Printf("First chapter length: %d characters\n", len(firstChapterHTML))
	
	// Check navigation
	hasNext := HasNextChapterByIndex(0)
	fmt.Printf("Has next chapter: %v\n", hasNext)
	
	// Get reading progress
	progress := GetReadingProgress(0)
	fmt.Printf("Reading progress: %.1f%%\n", progress)
	
	// Close the book
	CloseBook()
}
