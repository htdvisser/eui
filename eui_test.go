package eui_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"go.htdvisser.nl/eui"
)

func Example() {
	fmt.Println("Opening DB")
	db, err := bolt.Open(filepath.Join(os.Getenv("TMPDIR"), "eui-registrations.db"), 0600, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		fmt.Println("Closing DB")
		if err := db.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	fmt.Println("Initializing DB")
	err = eui.InitializeDB(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Fetching Registrations")
	registrations, err := eui.Fetch(eui.Registries[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Writing Registrations to DB")
	err = eui.WriteToDB(db, registrations)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Searching for 70B3D57ED000AB69")
	res, err := eui.FindInDB(db, []byte("70B3D57ED000AB69"))
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, res := range res {
		fmt.Printf("Result %d: %s\n", i+1, res.Name)
	}

	// Output:
	// Opening DB
	// Initializing DB
	// Fetching Registrations
	// Writing Registrations to DB
	// Searching for 70B3D57ED000AB69
	// Result 1: The Things Network Foundation
	// Closing DB
}
