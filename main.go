package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func exitHandler() error {
	fmt.Println("[SMDB] Exiting database server...")
	os.Exit(0)
	return nil
}

func insertHandler() error {
	fmt.Println("[SMDB] Inserting row...")
	return nil
}

func selectHandler() error {
	fmt.Println("[SMDB] Selecting row...")
	return nil
}

func main() {
	// Map of meta commands to their respective functions
	metaFuncMap := map[string]func() error{
		".exit": exitHandler,
	}

	stmtFuncMap := map[string]func() error{
		"insert": insertHandler,
		"select": selectHandler,
	}

	fmt.Println("[SMDB] Starting database server...")
	for {
		input, err := getInput()
		if err != nil {
			log.Fatal(err)
		}

		if input[0] == '.' {
			if metaFuncMap[input] != nil {
				err := metaFuncMap[input]()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("[SMDB] Unrecognized command \"%s\"\n", input)
			}
		} else {
			start := input[:6]
			fmt.Println(start)
			if stmtFuncMap[start] != nil {
				err := stmtFuncMap[start]()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("[SMDB] Unrecognized command \"%s\"\n", input)
			}
		}
	}
}

func getInput() (string, error) {
	fmt.Print("[SMDB] ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line[:len(line)-1], nil
}
