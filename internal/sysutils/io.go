package sysutils

import "fmt"

const (
	// cancelKey defines the key that cancels the confirmation.
	cancelKey = "n"
)

// Confirm prompts the user with a confirmation message and waits for input.
//
// The function displays the given message and asks the user to confirm by pressing any key
// to proceed or cancelKey to cancel. It reads input from standard input.
//
// Parameters:
//   - msg: The message to display to the user.
//
// Returns:
//   - true if the user presses any key other than cancelKey.
//   - false if the user presses cancelKey.
func Confirm(msg string) bool {
	fmt.Printf("\n%s (press n to cancel, or any key to proceed): ", msg)

	input := ""
	fmt.Scanln(&input)
	if input == cancelKey {
		return false
	}

	return true
}
