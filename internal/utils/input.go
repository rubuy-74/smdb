package utils

import (
	"bufio"
	"fmt"
	"os"
)

func GetInput() (string, error) {
	fmt.Print("[SMDB] > ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line[:len(line)-1], nil
}

func ToString(bytes []byte) string {
	return string(bytes[:])
}
