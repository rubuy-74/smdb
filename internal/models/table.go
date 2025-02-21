package models

const MAX_TABLE_SIZE = 100

type Table struct {
	Rows []string
}

var TableInstance = Table{}
