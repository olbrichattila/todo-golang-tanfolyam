package helpers

import (
	"bufio"
	"os"
)

func GetLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Remove newline characters
	input = input[:len(input)-1]

	return input, nil
}
