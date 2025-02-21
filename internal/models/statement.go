package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Statement struct {
	StmtType  string
	InsertRow Row
}

func PrepareStatement(input string) (*Statement, error) {
	stmtHead := input[:6]
	if stmtHead != "insert" && stmtHead != "select" {
		return &Statement{}, fmt.Errorf("unrecognized command \"%s\"", stmtHead)
	}
	stmt := &Statement{StmtType: stmtHead}
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
		stmt.InsertRow.ID = id
		stmt.InsertRow.Username = values[1]
		stmt.InsertRow.Email = values[2]
	}

	return stmt, nil
}
