package epubreader

import (
	"encoding/json"
	"fmt"
)

// Mobile-friendly wrapper functions for gomobile compatibility
// These functions use basic types (string, int, []byte) that work well with mobile bindings

// Global reader instance - mobile apps will work with one book at a time
var globalReader *Reader

// OpenEpubFromBytes opens an EPUB from byte array (mobile-compatible)
func OpenEpubFromBytes(data []byte) error {
	reader, err := OpenEpub(data)
	if err != nil {
		return err
	}
	globalReader = reader
	return nil
}

// GetBookInfo returns book metadata as JSON string
func GetBookInfo() string {
	if globalReader == nil {
		return `{"error": "no book loaded"}`
	}
	
	metadata := globalReader.GetMetadata()
	infoJSON, _ := json.Marshal(map[string]interface{}{
		"title":       metadata.Title,
		"author":      metadata.Author,
		"description": metadata.Description,
		"publisher":   metadata.Publisher,
		"language":    metadata.Language,
		"identifier":  metadata.Identifier,
		"chapterCount": globalReader.GetChapterCount(),
	})
	
	return string(infoJSON)
}

// GetTableOfContents returns TOC as JSON string
func GetTableOfContents() string {
	if globalReader == nil {
		return `{"error": "no book loaded"}`
	}
	
	chapters := globalReader.GetTOC()
	tocJSON, _ := json.Marshal(chapters)
	return string(tocJSON)
}

// GetChapterHTMLByIndex returns chapter HTML content by index
func GetChapterHTMLByIndex(index int) string {
	if globalReader == nil {
		return "<html><body><h1>Error: No book loaded</h1></body></html>"
	}
	
	html, err := globalReader.GetChapterHTML(index)
	if err != nil {
		return fmt.Sprintf("<html><body><h1>Error: %s</h1></body></html>", err.Error())
	}
	
	return html
}

// GetNextChapterHTMLByIndex returns next chapter HTML content
func GetNextChapterHTMLByIndex(currentIndex int) string {
	if globalReader == nil {
		return "<html><body><h1>Error: No book loaded</h1></body></html>"
	}
	
	html, err := globalReader.GetNextChapterHTML(currentIndex)
	if err != nil {
		return fmt.Sprintf("<html><body><h1>Error: %s</h1></body></html>", err.Error())
	}
	
	return html
}

// GetPrevChapterHTMLByIndex returns previous chapter HTML content
func GetPrevChapterHTMLByIndex(currentIndex int) string {
	if globalReader == nil {
		return "<html><body><h1>Error: No book loaded</h1></body></html>"
	}
	
	html, err := globalReader.GetPrevChapterHTML(currentIndex)
	if err != nil {
		return fmt.Sprintf("<html><body><h1>Error: %s</h1></body></html>", err.Error())
	}
	
	return html
}

// HasNextChapterByIndex checks if there's a next chapter
func HasNextChapterByIndex(currentIndex int) bool {
	if globalReader == nil {
		return false
	}
	return globalReader.HasNextChapter(currentIndex)
}

// HasPrevChapterByIndex checks if there's a previous chapter
func HasPrevChapterByIndex(currentIndex int) bool {
	if globalReader == nil {
		return false
	}
	return globalReader.HasPrevChapter(currentIndex)
}

// GetChapterTitleByIndex returns chapter title by index
func GetChapterTitleByIndex(index int) string {
	if globalReader == nil {
		return "Error: No book loaded"
	}
	return globalReader.GetChapterTitle(index)
}

// GetReadingProgress returns reading progress as percentage (0-100)
func GetReadingProgress(currentChapter int) float64 {
	if globalReader == nil {
		return 0.0
	}
	return globalReader.GetBookProgress(currentChapter)
}

// GetChapterCount returns total number of chapters
func GetChapterCount() int {
	if globalReader == nil {
		return 0
	}
	return globalReader.GetChapterCount()
}

// CloseBook closes the current book and frees memory
func CloseBook() {
	globalReader = nil
}
