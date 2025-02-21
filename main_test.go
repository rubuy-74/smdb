/* package main

import (
	"fmt"
	"testing"
)

func TestSerializeRow1(t *testing.T) {
	// test SerializeRow
	testRow := Row{id: 1, username: "test", email: "test@example"}
	serializedRow := SerializeRow(testRow)
	expected := "1;test;test@example"
	if serializedRow != expected {
		t.Errorf("expected %s but got %s", expected, serializedRow)
	}
}

func TestSerializeRow2(t *testing.T) {
	// test SerializeRow
	testRow := Row{username: "test2", email: "test2@example"}
	serializedRow := SerializeRow(testRow)
	expected := "0;test2;test2@example"
	if serializedRow != expected {
		t.Errorf("expected %s but got %s", expected, serializedRow)
	}
}

func TestDeserializeRow1(t *testing.T) {
	// test DeserializeRow
	testRow := "1;test;test@example"
	deserializedRow := DeserializeRow(testRow)
	expected := Row{id: 1, username: "test", email: "test@example"}
	if deserializedRow != expected {
		t.Errorf("expected %v but got %v", expected, deserializedRow)
	}
}

func TestDeserializeRow2(t *testing.T) {
	// test DeserializeRow
	testRow := "0;test2;test2@example"
	deserializedRow := DeserializeRow(testRow)
	expected := Row{username: "test2", email: "test2@example"}
	if deserializedRow != expected {
		t.Errorf("expected %v but got %v", expected, deserializedRow)
	}
}

func TestInsertHandler(t *testing.T) {
	// test InsertHandler
	stmt := &Statement{stmtType: "insert", insertRow: Row{id: 1, username: "test", email: "test@example"}}
	err := InsertHandler(stmt)
	if err != nil {
		t.Errorf("expected nil but got %v", err)
	}
}

func TestInsertHandlerFullTable(t *testing.T) {
	// test InsertHandler
	stmt := &Statement{stmtType: "insert", insertRow: Row{id: 1, username: "test", email: "test@example"}}
	for i := 0; i < 100; i++ {
		_ = InsertHandler(stmt)
	}
	err := InsertHandler(stmt)
	expected_error := fmt.Errorf("table is full")
	if err.Error() != expected_error.Error() {
		t.Errorf("expected (%v) but got (%v)", expected_error, err)
	}
}

func TestSelectHandler(t *testing.T) {
	// test SelectHandler
	stmt := &Statement{stmtType: "select"}
	err := SelectHandler(stmt)
	if err != nil {
		t.Errorf("expected nil but got %v", err)
	}
}

func TestPrepareStatementInsert(t *testing.T) {
	// test prepareStatement
	input := "insert 1 test test@example"
	stmt, err := prepareStatement(input)
	if err != nil {
		t.Errorf("expected nil but got %v", err)
	}
	expected := &Statement{stmtType: "insert", insertRow: Row{id: 1, username: "test", email: "test@example"}}
	if stmt.stmtType != expected.stmtType || stmt.insertRow != expected.insertRow {
		t.Errorf("expected %v but got %v", expected, stmt)
	}
}

func TestPrepareStatementInsertIncorrect(t *testing.T) {
	// test prepareStatement
	input := "insert"
	_, err := prepareStatement(input)
	expected := fmt.Errorf("invalid insert statement")
	if err.Error() != expected.Error() {
		t.Errorf("expected (%v) but got (%v)", expected, err)
	}
}

func TestPrepareStatementSelect(t *testing.T) {
	// test prepareStatement
	input := "select"
	stmt, err := prepareStatement(input)
	if err != nil {
		t.Errorf("expected nil but got %v", err)
	}
	expected := &Statement{stmtType: "select"}
	if stmt.stmtType != expected.stmtType {
		t.Errorf("expected (%v) but got (%v)", expected, stmt)
	}
}

func TestPrepareStatementIncorrect(t *testing.T) {
	// test prepareStatement
	input := "incorrect"
	_, err := prepareStatement(input)
	expected := fmt.Errorf("unrecognized command \"incorrect\"")
	if err.Error() != expected.Error() {
		t.Errorf("expected (%v) but got (%v)", expected, err)
	}
}
*/