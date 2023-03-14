package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

func init() {
	rootCMD.AddCommand(doCMD)
}

var doCMD = &cobra.Command{
	Use:   "do",
	Args:  checkArgs,
	Short: "Marks a TODO as done",
	Run:   doFunc,

	Example: "todo do 1",
}

func doFunc(cmd *cobra.Command, args []string) {
	fmt.Println("[debug] using db:", dbPath)
	db, err := bolt.Open(dbPath, 0o666, nil)
	if err != nil {
		log.Fatalf("error opening db: %s\n", err)
	}
	defer db.Close()

	todoID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("error getting ID of TODO %s: %s\n", args[0], err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		todoBucket := tx.Bucket([]byte("TODO"))

		todoStr := todoBucket.Get(itob(todoID))
		if todoStr == nil {
			log.Printf("the TODO %d does not exist!\n", todoID)
			os.Exit(1)
		}
		todoBucket.Delete(itob(todoID))
		if err != nil {
		}

		fmt.Printf("completing TODO %d: %s\n", todoID, string(todoStr))

		return nil
	})
	if err != nil {
		log.Fatalf("error getting TODOs from db: %s\n", err)
	}
}
