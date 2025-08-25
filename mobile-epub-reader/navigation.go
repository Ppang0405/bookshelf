package epubreader

import "fmt"

// Navigation functions for mobile apps

// GetNextChapterHTML returns the HTML content of the next chapter
func (e *Reader) GetNextChapterHTML(currentIndex int) (string, error) {
	nextIndex := currentIndex + 1
	if nextIndex >= len(e.chapters) {
		return "", fmt.Errorf("no next chapter available")
	}
	return e.GetChapterHTML(nextIndex)
}

// GetPrevChapterHTML returns the HTML content of the previous chapter
func (e *Reader) GetPrevChapterHTML(currentIndex int) (string, error) {
	prevIndex := currentIndex - 1
	if prevIndex < 0 {
		return "", fmt.Errorf("no previous chapter available")
	}
	return e.GetChapterHTML(prevIndex)
}

// HasNextChapter checks if there's a next chapter
func (e *Reader) HasNextChapter(currentIndex int) bool {
	return currentIndex+1 < len(e.chapters)
}

// HasPrevChapter checks if there's a previous chapter
func (e *Reader) HasPrevChapter(currentIndex int) bool {
	return currentIndex > 0
}

// GetChapterTitle returns the title of a specific chapter
func (e *Reader) GetChapterTitle(chapterIndex int) string {
	if chapterIndex < 0 || chapterIndex >= len(e.chapters) {
		return ""
	}
	return e.chapters[chapterIndex].Title
}

// FindChapterByTitle searches for a chapter by title (case-insensitive)
func (e *Reader) FindChapterByTitle(title string) int {
	for i, chapter := range e.chapters {
		if chapter.Title == title {
			return i
		}
	}
	return -1 // Not found
}

// GetBookProgress returns reading progress as percentage (0-100)
func (e *Reader) GetBookProgress(currentChapter int) float64 {
	if len(e.chapters) == 0 {
		return 0.0
	}
	progress := float64(currentChapter) / float64(len(e.chapters)) * 100
	if progress > 100 {
		progress = 100
	}
	return progress
}
