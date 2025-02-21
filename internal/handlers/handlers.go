package handlers

import (
	"fmt"
	"os"

	"github.com/rubuy-74/smDB/internal/models"
)

func ExitHandler() error {
	fmt.Println("[SMDB] Exiting database server...")
	os.Exit(0)
	return nil
}

func InsertHandler(stmt *models.Statement) error {
	serializedRow := models.SerializeRow(stmt.InsertRow)
	if len(models.TableInstance.Rows) < models.MAX_TABLE_SIZE {
		models.TableInstance.Rows = append(models.TableInstance.Rows, serializedRow)
	} else {
		return fmt.Errorf("table is full")
	}
	return nil
}

func SelectHandler(stmt *models.Statement) error {
	for _, row := range models.TableInstance.Rows {
		deserializedRow := models.DeserializeRow(row)
		fmt.Printf("(%d,%s,%s)\n", deserializedRow.ID, deserializedRow.Username, deserializedRow.Email)
	}
	return nil
}
