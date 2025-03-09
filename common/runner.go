package common



type runnerRef struct {
	Llama          interface{} // Exported field
	SessionDuration int         // Exported field
	NumParallel    int         // Exported field
	RefCount       int         // Exported field
	Gpus           []string    // Exported field
	Loading        bool        // Exported field
	// Removed direct reference to api.Model to avoid circular dependency
}
