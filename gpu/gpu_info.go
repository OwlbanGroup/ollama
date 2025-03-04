package gpu

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lgpu_info
#include "gpu_info_darwin.h"
*/
import "C"
import (
	"log"
)

// GetCPUInfo retrieves information about the CPU.
func GetCPUInfo() (uint64, uint64) {
	totalMemory := getPhysicalMemory() // Example value in MB
	freeMemory := getFreeMemory()       // Example value in MB
	log.Printf("CPU Info - Total Memory: %d, Free Memory: %d", totalMemory, freeMemory)
	return totalMemory, freeMemory // Adjusted to return only two values
}

// GetGPUInfo retrieves information about the GPU.
func GetGPUInfo() (uint64, uint64, uint64) {
	vram := getRecommendedMaxVRAM() // Get recommended max VRAM
	physicalMemory := getPhysicalMemory() // Example value in MB
	freeMemory := getFreeMemory()       // Example value in MB
	log.Printf("GPU Info - VRAM: %d, Physical Memory: %d, Free Memory: %d", vram, physicalMemory, freeMemory)
	return vram, physicalMemory, freeMemory
}

// GetVisibleDevicesEnv returns the environment variable for visible devices.
func GetVisibleDevicesEnv() (string, string) {
	return "VISIBLE_DEVICES", "0" // Example values
}
