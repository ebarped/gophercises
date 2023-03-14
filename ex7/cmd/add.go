package cmd

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

func init() {
	rootCMD.AddCommand(addCMD)
}

var addCMD = &cobra.Command{
	Use:       "add",
	Args:      checkArgs,
	ValidArgs: []string{"Buy the newspaper"},
	Short:     "Adds a new TODO",
	Run:       addFunc,
	Example:   "todo add \"buy newspaper\"",
}

func addFunc(cmd *cobra.Command, args []string) {
	fmt.Println("[debug] using db:", dbPath)
	db, err := bolt.Open(dbPath, 0o666, nil)
	if err != nil {
		log.Fatalf("error opening db: %s\n", err)
	}
	defer db.Close()

	todoStr := args[0]

	err = db.Update(func(tx *bolt.Tx) error {
		todoBucket := tx.Bucket([]byte("TODO"))
		id, err := todoBucket.NextSequence()
		if err != nil {
			log.Fatalf("error generating ID for this TODO: %s\n", err)
		}

		err = todoBucket.Put(itob(int(id)), []byte(todoStr))
		if err != nil {
			log.Fatalf("error putting TODO %s: %s\n", todoStr, err)
		}
		fmt.Printf("adding TODO: %q\n", todoStr)

		return nil
	})
	if err != nil {
		log.Fatalf("error getting TODOs from db: %s\n", err)
	}
}

// checkArgs perform various checks over the arguments
func checkArgs(cmd *cobra.Command, args []string) error {
	err := cobra.ExactArgs(1)(cmd, args)
	if err != nil {
		fmt.Println("[debug] args passed:", args)
		return err
	}
	return nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
