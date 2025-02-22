package models

import "fmt"

const MAX_TABLE_SIZE = 100
const MAX_PAGE_SIZE = 4096
const ROW_SIZE = 296
const MAX_NUM_ROW_PAGE = 4096 / 296

type Page struct {
	numRows int
	rows    [MAX_PAGE_SIZE]byte
}

func (p *Page) AddRow(row Row) error {
	if p.numRows >= 13 {
		return fmt.Errorf("max num rows in page")
	}
	copy(p.rows[(ROW_SIZE*p.numRows):], SerializeRow(row))
	p.numRows += 1
	return nil
}

func (p *Page) GetRow(rowId int) Row {
	rowOffset := rowId * 296
	rowBytes := p.rows[rowOffset : rowOffset+296]
	row := DeserializeRow(rowBytes)
	return row
}

// var TableInstance = Table{}
