// Package rethink provides ...
package rethink

import (
	"log"
)

// LogInfo log information
func LogInfo(message string, args ...interface{}) {
	log.Printf(message, args)
}

// LogError log error message
func LogError(message string, args ...interface{}) {
	log.Printf(message, args)
}

// LogDebug log debug info
func LogDebug(message string, args ...interface{}) {
	log.Printf(message, args)
}
