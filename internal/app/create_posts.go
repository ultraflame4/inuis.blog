package app

// Generates a title and summary for the post based on the content.
// Returns (title, summary)
func GeneratePostTitleAndSummary(content string) (string, string) {
	title := ""
	body := ""
	// Remove the first line of the content. (Treated as title)
	if i := strings.IndexByte(body, '\n'); i >= 0 {
		title = content[:i]
		body = content[i+1:]
	} else {
		// If i == -1. Means no newlines. The entire content is the first line.
		title = strings.Clone(content)
		// Also means that body has no content
		body = ""
	}
	title = strings.TrimPrefix(title, "#")
	title = strings.TrimSpace(title)

	// Generate summary from body
	summary_len := 100
	summary_body := strings.ReplaceAll(body, "\n", "")
	summary := summary_body[0:min(len(summary_body), summary_len)]
	summary = strings.TrimSpace(summary)
	return title, summary
}

// func CreatePost(content string, db *gorm.DB, ctx context.Context) (models.Post, error) {

// }
