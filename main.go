package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const MAX_TABLE_SIZE = 100

var metaFuncMap map[string]func() error
var table Table
var stmtFuncMap = map[string]func(*Statement) error{
	"insert": InsertHandler,
	"select": SelectHandler,
}

type Table struct {
	rows []string
}

type Row struct {
	id       int
	username string
	email    string
}

type Statement struct {
	stmtType  string
	insertRow Row
}

func GetInput() (string, error) {
	fmt.Print("[SMDB] > ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line[:len(line)-1], nil
}

func ExitHandler() error {
	fmt.Println("[SMDB] Exiting database server...")
	os.Exit(0)
	return nil
}

func InsertHandler(stmt *Statement) error {
	serializedRow := SerializeRow(stmt.insertRow)
	if len(table.rows) < MAX_TABLE_SIZE {
		table.rows = append(table.rows, serializedRow)
	} else {
		return fmt.Errorf("table is full")
	}
	return nil
}

func SelectHandler(stmt *Statement) error {
	for _, row := range table.rows {
		deserializedRow := DeserializeRow(row)
		fmt.Printf("(%d,%s,%s)\n", deserializedRow.id, deserializedRow.username, deserializedRow.email)
	}
	return nil
}

func SerializeRow(row Row) string {
	serializedRow := fmt.Sprintf("%d;%s;%s", row.id, row.username, row.email)
	return serializedRow
}

func DeserializeRow(serializedRow string) Row {
	row := Row{}
	values := strings.Split(serializedRow, ";")
	row.id, _ = strconv.Atoi(values[0])
	row.username = values[1]
	row.email = values[2]
	return row
}

func main() {
	metaFuncMap = map[string]func() error{
		".exit": ExitHandler,
	}

	fmt.Println("[SMDB] Starting database server...")
	for {
		input, err := GetInput()
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
			stmt, err := prepareStatement(input)
			if err != nil {
				fmt.Printf("[SMDB] Preparing command - \"%s\"\n", err)
			} else {
				err = stmtFuncMap[stmt.stmtType](stmt)
				if err != nil {
					fmt.Printf("[SMDB] Execution error - \"%s\"\n", err)
				} else {
					fmt.Println("[SMDB] Command executed successfully")
				}
			}
		}
	}
}

// insert 1 cstack foo@bar.com
func prepareStatement(input string) (*Statement, error) {
	stmtHead := strings.Split(input, " ")[0]
	if stmtFuncMap[stmtHead] == nil {
		return &Statement{}, fmt.Errorf("unrecognized command \"%s\"", stmtHead)
	}
	stmt := &Statement{stmtType: stmtHead}
	if stmtHead == "insert" {
		// parse insert statement
		values := strings.Split(input, " ")[1:]
		if len(values) < 3 {
			return &Statement{}, fmt.Errorf("invalid insert statement")
		}

		id, err := strconv.Atoi(values[0])
		if err != nil {
			return &Statement{}, fmt.Errorf("invalid id value")
		}
		stmt.insertRow.id = id
		stmt.insertRow.username = values[1]
		stmt.insertRow.email = values[2]
	}

	return stmt, nil
}
