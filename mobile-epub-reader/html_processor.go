package epubreader

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// GetChapterHTML returns the HTML content of a specific chapter
func (e *Reader) GetChapterHTML(chapterIndex int) (string, error) {
	if chapterIndex < 0 || chapterIndex >= len(e.chapters) {
		return "", fmt.Errorf("chapter index %d out of range", chapterIndex)
	}

	// Get chapter content using the epub reader
	content, err := e.GetChapterContent(chapterIndex)
	if err != nil {
		return "", fmt.Errorf("failed to read chapter content: %w", err)
	}

	// Clean and process the HTML
	cleanHTML, err := e.processHTML(content)
	if err != nil {
		return "", fmt.Errorf("failed to process HTML: %w", err)
	}

	return cleanHTML, nil
}

// processHTML cleans and prepares HTML for mobile display
func (e *Reader) processHTML(rawHTML string) (string, error) {
	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return "", err
	}

	// Extract body content
	bodyContent := e.extractBodyContent(doc)
	
	// Clean the HTML
	cleanedContent := e.cleanHTML(bodyContent)
	
	// Wrap in mobile-friendly HTML
	mobileHTML := e.wrapForMobile(cleanedContent)
	
	return mobileHTML, nil
}

// extractBodyContent extracts content from HTML body
func (e *Reader) extractBodyContent(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "body" {
		return e.renderNode(node)
	}
	
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if content := e.extractBodyContent(child); content != "" {
			return content
		}
	}
	
	// If no body found, render the entire document
	return e.renderNode(node)
}

// renderNode converts an HTML node to string
func (e *Reader) renderNode(node *html.Node) string {
	var buf bytes.Buffer
	html.Render(&buf, node)
	return buf.String()
}

// cleanHTML removes unwanted elements and attributes
func (e *Reader) cleanHTML(htmlContent string) string {
	// Remove script tags
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	htmlContent = scriptRegex.ReplaceAllString(htmlContent, "")
	
	// Remove style tags (we'll add our own)
	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	htmlContent = styleRegex.ReplaceAllString(htmlContent, "")
	
	// Remove comments
	commentRegex := regexp.MustCompile(`<!--.*?-->`)
	htmlContent = commentRegex.ReplaceAllString(htmlContent, "")
	
	// Clean up whitespace
	whitespaceRegex := regexp.MustCompile(`\s+`)
	htmlContent = whitespaceRegex.ReplaceAllString(htmlContent, " ")
	
	return strings.TrimSpace(htmlContent)
}

// wrapForMobile wraps content in mobile-friendly HTML
func (e *Reader) wrapForMobile(content string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            margin: 16px;
            background-color: #ffffff;
            color: #333333;
        }
        
        p {
            margin: 0 0 1em 0;
            text-align: justify;
        }
        
        h1, h2, h3, h4, h5, h6 {
            margin: 1.5em 0 0.5em 0;
            line-height: 1.2;
        }
        
        img {
            max-width: 100%%;
            height: auto;
        }
        
        blockquote {
            margin: 1em 20px;
            padding: 10px 20px;
            background-color: #f5f5f5;
            border-left: 4px solid #ddd;
        }
        
        /* Dark mode support */
        @media (prefers-color-scheme: dark) {
            body {
                background-color: #1a1a1a;
                color: #e0e0e0;
            }
            blockquote {
                background-color: #333;
                border-left-color: #555;
            }
        }
    </style>
</head>
<body>
    %s
</body>
</html>`, e.metadata.Title, content)
}