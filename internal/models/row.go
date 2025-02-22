package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Row struct {
	id       int64
	username [32]byte
	email    [256]byte
}

func CreateRow(id int64, username string, email string) (Row, error) {
	if len(username) > 32 {
		return Row{}, fmt.Errorf("invalid username size")
	}
	if len(email) > 256 {
		return Row{}, fmt.Errorf("invalid email size")
	}
	row := Row{}

	row.SetID(id)
	row.SetUsername(username)
	row.SetEmail(email)

	return row, nil
}

func (r *Row) SetUsername(username string) {
	copy(r.username[:], string(username))
}

func (r *Row) SetEmail(email string) {
	copy(r.email[:], string(email))
}

func (r *Row) SetID(id int64) {
	r.id = id
}

// returns [295]byte array with the information on the row
// id: 8 bytes
// username: 32 bytes
// email: 256 bytes
func SerializeRow(row Row) []byte {
	var rowSerialized bytes.Buffer
	binary.Write(&rowSerialized, binary.LittleEndian, row.id)
	rowSerialized.Write(row.username[:])
	rowSerialized.Write(row.email[:])
	return rowSerialized.Bytes()
}

// returns the original Row given the serialized array
func DeserializeRow(serializedRow []byte) Row {
	var id int64
	var username [32]byte
	var email [256]byte

	idBuffer := bytes.NewReader(serializedRow[:8])
	binary.Read(idBuffer, binary.LittleEndian, &id)
	copy(username[:], serializedRow[8:40])
	copy(email[:], serializedRow[40:])

	row := Row{id: id, username: username, email: email}
	return row
}

func PrintRow(row Row) {
	id := row.id
	username := string(bytes.Trim(row.username[:], "\x00"))
	email := string(bytes.Trim(row.email[:], "\x00"))
	fmt.Printf("(%d,%s,%s)\n", id, username, email)
}
