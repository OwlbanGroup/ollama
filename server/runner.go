package server

type runnerRef struct {
	llama          interface{} // Replace with the actual type
	sessionDuration int
	numParallel    int
	refCount       int
	gpus           []string // Adjust type as necessary
	loading        bool
	model          *Model // Replace with the actual type
}
