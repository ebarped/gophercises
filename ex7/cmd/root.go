package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const (
	CONFIG_FOLDER = ".todo"
	DB_NAME       = "bolt.db"
)

var (
	dbParentPath string // path of the parent dir of the db file
	dbPath       string // path to the db file
)

var rootCMD = &cobra.Command{
	Use:              "todo",
	Short:            "tool to manage TODOs (add, list and complete your TODOs)",
	PersistentPreRun: initConfig,
	Run:              rootFunc,
}

// execute common initial steps
func initConfig(cmd *cobra.Command, args []string) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error: error getting the home directory: %s\n", err)
	}

	dbParentPath = userHome + "/" + CONFIG_FOLDER
	dbPath = dbParentPath + "/" + DB_NAME

	setupDB()
}

func rootFunc(cmd *cobra.Command, args []string) {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	listFunc(cmd.Root(), []string{})
}

// Execute adds all child commands to the root command, and sets flags
func Execute() {
	if err := rootCMD.Execute(); err != nil {
		log.Fatalf("error on startup: %s\n", err)
	}
}

func setupDB() {
	fmt.Println("[debug] Executing initial DB setup using db:", dbPath)

	err := os.MkdirAll(dbParentPath, os.ModePerm)
	if err != nil {
		log.Fatalf("error creating path to store db: %s\n", err)
	}

	fmt.Println("[debug] using db:", dbPath)
	db, err := bolt.Open(dbPath, 0o666, nil)
	if err != nil {
		log.Fatalf("error opening db: %s\n", err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("TODO"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("could not set up buckets, %s\n", err)
	}
	fmt.Println("[debug] DB Setup Done")
}
