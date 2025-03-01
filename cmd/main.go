package main

import (
	"log"

	"github.com/rubuy-74/smDB/internal/b3"
)

/* var metaFuncMap map[string]func() error
var stmtFuncMap = map[string]func(*models.Statement) error{
	"insert": handlers.InsertHandler,
	"select": handlers.SelectHandler,
} */

func main() {
}

func updateKV() {
	node := createdFirstNode()
	key := []byte("k3")
	newVal := []byte("updated")
	new, err := b3.ChangeKVPair(node, key, newVal)
	if err != nil {
		log.Fatal(err)
	}
	new.PrintAllKV()
}

func createdFirstNode() b3.BNode {
	node := b3.BNode(make([]byte, b3.BTREE_PAGE_SIZE))
	node.SetHeader(b3.BNODE_LEAF, 3)
	b3.NodeAppendKV(node, 0, 0, []byte("k1"), []byte("hiaaah"))
	b3.NodeAppendKV(node, 1, 0, []byte("k3"), []byte("hello"))
	b3.NodeAppendKV(node, 2, 0, []byte("k4"), []byte("fucku"))

	return node
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
