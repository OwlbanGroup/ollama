package llm

// projectorMemoryRequirements defines the memory requirements for the projector.
var projectorMemoryRequirements = map[string]int{
    "low":    512,  // 512 MB for low memory requirements
    "medium": 1024, // 1 GB for medium memory requirements
    "high":   2048, // 2 GB for high memory requirements
}
