package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("[SMDB] Starting database server...")
	for {
		fmt.Print("[SMDB] ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatal(err)
		}

		if input == ".exit" {
			fmt.Println("[SMDB] Exiting database server...")
			return
		} else {
			fmt.Println("[SMDB] {ERROR} Command not found {ERROR}")
		}
	}
}
