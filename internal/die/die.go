// Package die provides custom panic function for error messages with execution ending
package die

import (
	"fmt"
	"os"
)

func Die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
