package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
	"go.htdvisser.nl/eui"
)

const defaultDatabase = "$TMPDIR/eui-registrations.db"

var (
	database   = flag.String("db", defaultDatabase, "Database file")
	initialize = flag.Bool("initialize", false, "Initialize the registrations database")
	update     = flag.Bool("update", false, "Update the registrations database")
)

func colorize(msg string) string {
	if os.Getenv("COLORTERM") != "" {
		return strings.Replace(msg, "%s", "\033[34m%6s\033[0m", -1)
	}
	return msg
}

func main() {
	flag.Parse()

	dbLocation := *database
	if dbLocation == defaultDatabase {
		tmpDir := os.Getenv("TMPDIR")
		if tmpDir == "" {
			tmpDir = "/tmp"
		}
		dbLocation = filepath.Join(tmpDir, "eui-registrations.db")
	}

	db, err := bolt.Open(dbLocation, 0600, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
			return
		}
	}()

	if *initialize {
		log.Println("Initializing database...")
		err := eui.InitializeDB(db)
		if err != nil {
			log.Print(err)
			return
		}
	}

	if *initialize || *update {
		for _, registry := range eui.Registries {
			log.Printf(colorize("Fetching %s..."), registry)
			registrations, err := eui.Fetch(registry)
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("Updating database...")
			err = eui.WriteToDB(db, registrations)
			if err != nil {
				log.Print(err)
				return
			}
			log.Println("Updated database")
		}
	}

	for _, search := range flag.Args() {
		search = strings.ToUpper(search)
		regs, err := eui.FindInDB(db, []byte(search))
		if err != nil {
			log.Println(err)
			continue
		}
		for _, reg := range regs {
			log.Printf(colorize("EUI %s belongs to %s (prefix %s)"), search, reg.Name, reg.Prefix)
		}
	}
}
