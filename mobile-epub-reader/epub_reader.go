package epubreader

import (
	"bytes"
	"fmt"
	"io"

	"github.com/taylorskalyo/goreader/epub"
)

// BookMetadata contains EPUB metadata
type BookMetadata struct {
	Title       string
	Author      string
	Description string
	Publisher   string
	Language    string
	Identifier  string
}

// Chapter represents a chapter in the EPUB
type Chapter struct {
	Title    string
	Index    int
	ID       string
	Filepath string
}

// Reader handles EPUB parsing and content extraction
type Reader struct {
	reader   *epub.Reader
	chapters []Chapter
	metadata BookMetadata
	zipData  []byte
}

// OpenEpub creates a new Reader from EPUB data
func OpenEpub(data []byte) (*Reader, error) {
	// Create a reader from the byte slice
	reader := bytes.NewReader(data)
	
	// Create epub reader
	epubReader, err := epub.NewReader(reader, int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create epub reader: %w", err)
	}

	e := &Reader{
		reader:  epubReader,
		zipData: data,
	}

	// Extract metadata
	if err := e.extractMetadata(); err != nil {
		return nil, fmt.Errorf("failed to extract metadata: %w", err)
	}

	// Extract chapters
	if err := e.extractChapters(); err != nil {
		return nil, fmt.Errorf("failed to extract chapters: %w", err)
	}

	return e, nil
}

// extractMetadata extracts book metadata from the EPUB
func (e *Reader) extractMetadata() error {
	if len(e.reader.Rootfiles) == 0 {
		return fmt.Errorf("no rootfiles found")
	}

	rootfile := e.reader.Rootfiles[0]
	metadata := rootfile.Metadata

	e.metadata = BookMetadata{
		Title:       metadata.Title,
		Author:      metadata.Creator,
		Description: metadata.Description,
		Publisher:   metadata.Publisher,
		Language:    metadata.Language,
		Identifier:  metadata.Identifier,
	}
	
	// Set defaults if empty
	if e.metadata.Title == "" {
		e.metadata.Title = "Unknown Title"
	}
	if e.metadata.Author == "" {
		e.metadata.Author = "Unknown Author"
	}

	return nil
}

// extractChapters extracts chapter information from the EPUB
func (e *Reader) extractChapters() error {
	if len(e.reader.Rootfiles) == 0 {
		return fmt.Errorf("no rootfiles found")
	}

	rootfile := e.reader.Rootfiles[0]
	spine := rootfile.Spine
	manifest := make(map[string]*epub.Item)
	
	// Build manifest map
	for i := range rootfile.Manifest.Items {
		item := &rootfile.Manifest.Items[i]
		manifest[item.ID] = item
	}

	e.chapters = make([]Chapter, 0, len(spine.Itemrefs))

	for i, itemref := range spine.Itemrefs {
		// Get the manifest item
		manifestItem := manifest[itemref.IDREF]
		if manifestItem == nil {
			continue
		}

		chapter := Chapter{
			Index:    i,
			ID:       itemref.IDREF,
			Filepath: manifestItem.HREF,
			Title:    fmt.Sprintf("Chapter %d", i+1), // Default title
		}

		e.chapters = append(e.chapters, chapter)
	}

	return nil
}

// GetMetadata returns the book metadata
func (e *Reader) GetMetadata() BookMetadata {
	return e.metadata
}

// GetChapterCount returns the number of chapters
func (e *Reader) GetChapterCount() int {
	return len(e.chapters)
}

// GetTOC returns the table of contents
func (e *Reader) GetTOC() []Chapter {
	return e.chapters
}

// GetChapterContent returns the content of a specific chapter
func (e *Reader) GetChapterContent(chapterIndex int) (string, error) {
	if chapterIndex < 0 || chapterIndex >= len(e.chapters) {
		return "", fmt.Errorf("chapter index %d out of range", chapterIndex)
	}

	chapter := e.chapters[chapterIndex]
	
	// Find the item in the manifest
	if len(e.reader.Rootfiles) == 0 {
		return "", fmt.Errorf("no rootfiles found")
	}

	rootfile := e.reader.Rootfiles[0]
	var item *epub.Item
	
	for i := range rootfile.Manifest.Items {
		if rootfile.Manifest.Items[i].ID == chapter.ID {
			item = &rootfile.Manifest.Items[i]
			break
		}
	}
	
	if item == nil {
		return "", fmt.Errorf("item not found for chapter: %s", chapter.ID)
	}

	// Open the item
	rc, err := item.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open chapter content: %w", err)
	}
	defer rc.Close()

	// Read the content
	content, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("failed to read chapter content: %w", err)
	}

	return string(content), nil
}
