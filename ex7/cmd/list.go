package cmd

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

func init() {
	rootCMD.AddCommand(listCMD)
}

var listCMD = &cobra.Command{
	Use:   "list",
	Short: "Lists all TODOs stored",
	Run:   listFunc,
}

func listFunc(cmd *cobra.Command, args []string) {
	fmt.Println("[debug] using db:", dbPath)

	db, err := bolt.Open(dbPath, 0o666, nil)
	if err != nil {
		log.Fatalf("error opening db: %s\n", err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		todoBucket := tx.Bucket([]byte("TODO"))
		keysNumber := todoBucket.Stats().KeyN
		if keysNumber == 0 {
			fmt.Println("All your TODOs are completed!")
			os.Exit(0)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error getting TODOs from db: %s\n", err)
	}

	fmt.Println("--- TODO LIST ---")
	err = db.View(func(tx *bolt.Tx) error {
		todoBucket := tx.Bucket([]byte("TODO"))
		todoBucket.ForEach(func(k, v []byte) error {
			fmt.Printf("%d. %s\n", btoi(k), v)
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatalf("error getting TODOs from db: %s\n", err)
	}
}

// itob returns a uint64 big endian representation of b.
func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
