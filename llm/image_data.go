package llm

type ImageData struct {
	ID     string `json:"id"` // Adding ID field
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}
