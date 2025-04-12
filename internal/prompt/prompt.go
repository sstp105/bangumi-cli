package prompt

import (
	"bufio"
	"fmt"
	"os"
)

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
	_, _ = fmt.Scanln(&input)
	if input == cancelKey {
		return false
	}

	return true
}

// ReadUserInput prompts the user with a message and reads a line of input from the standard input (stdin).
// It returns the input as a string.
func ReadUserInput(msg string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("\n%s: ", msg)

	scanner.Scan()
	return scanner.Text()
}
