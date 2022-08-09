package version

import "fmt"

var falcoctlVersion = "development"

// Run executes the handler for the version command.
func Run() error {
	fmt.Printf("Version: %s\n", falcoctlVersion)
	return nil
}
