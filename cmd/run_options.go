package cmd

type Message struct {
	Role    string
	Content string
}



// runOptions holds the options for building a model file.
type runOptions struct {
	Model      string
	System     string
	Messages   []Message // Updated to use local Message struct
	Options    map[string]interface{}
	ParentModel string // New field added
}
