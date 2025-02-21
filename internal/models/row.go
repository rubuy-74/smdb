package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Row struct {
	ID       int
	Username string
	Email    string
}

func SerializeRow(row Row) string {
	serializedRow := fmt.Sprintf("%d;%s;%s", row.ID, row.Username, row.Email)
	return serializedRow
}

func DeserializeRow(serializedRow string) Row {
	row := Row{}
	values := strings.Split(serializedRow, ";")
	row.ID, _ = strconv.Atoi(values[0])
	row.Username = values[1]
	row.Email = values[2]
	return row
}
