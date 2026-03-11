package parser

type Document struct {
	ID      string   `json:"id"`
	Source  string   `json:"source"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	URL     string   `json:"url,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}
