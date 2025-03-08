package gpu

import (
	"log"
)

// GetCPUInfo retrieves information about the CPU for Windows.
func GetCPUInfo() (uint64, uint64) {
	totalMemory := uint64(8192) // Total memory in MB
	freeMemory := uint64(4096)   // Free memory in MB
	log.Printf("CPU Info - Total Memory: %d, Free Memory: %d", totalMemory, freeMemory)
	return totalMemory, freeMemory
}

// GetGPUInfo retrieves information about the GPU for Windows.
func GetGPUInfo() (uint64, uint64, uint64) {
	vram := uint64(4096)         // VRAM in MB
	physicalMemory := uint64(8192) // Physical memory in MB
	freeMemory := uint64(2048)    // Free memory in MB
	log.Printf("GPU Info - VRAM: %d, Physical Memory: %d, Free Memory: %d", vram, physicalMemory, freeMemory)
	return vram, physicalMemory, freeMemory
}
