package config

import "os"

// if true then we are in test mode.
func IsInTests() bool {
	return os.Getenv("IN_TESTS") == "1"
}
