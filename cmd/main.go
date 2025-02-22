package main

import (
	"github.com/rubuy-74/smDB/internal/models"
)

const MAX_TABLE_SIZE = 100

/* var metaFuncMap map[string]func() error
var stmtFuncMap = map[string]func(*models.Statement) error{
	"insert": handlers.InsertHandler,
	"select": handlers.SelectHandler,
} */

func main() {
	page := models.Page{}
	rowTest, _ := models.CreateRow(0, "a", "a")
	rowTest2, _ := models.CreateRow(1, "b", "b")
	rowTest3, _ := models.CreateRow(2, "c", "c")
	page.AddRow(rowTest)
	page.AddRow(rowTest2)
	page.AddRow(rowTest3)
	/* 	models.PrintRow(page.GetRow(0))
	   	models.PrintRow(page.GetRow(1))
	   	models.PrintRow(page.GetRow(2)) */

}

/* func main() {


		 	metaFuncMap = map[string]func() error{
				".exit": handlers.ExitHandler,
			}

			fmt.Println("[SMDB] Starting database server...")
			for {
				input, err := utils.GetInput()
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
					stmt, err := models.PrepareStatement(input)
					if err != nil {
						fmt.Printf("[SMDB] Preparing command - \"%s\"\n", err)
					} else {
						err = stmtFuncMap[stmt.StmtType](stmt)
						if err != nil {
							fmt.Printf("[SMDB] Execution error - \"%s\"\n", err)
						} else {
							fmt.Println("[SMDB] Command executed successfully")
						}
					}
				}
			}

}
*/
