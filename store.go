package eui

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/boltdb/bolt"
)

const bucketName = "Registrations"

var ErrNotInitialized = errors.New("Database not initialized")

func WriteToDB(db *bolt.DB, registrations []Registration) error {
	sort.Sort(Registrations(registrations))
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrNotInitialized
		}
		for _, reg := range registrations {
			if err := bucket.Put([]byte(strings.ToUpper(reg.Prefix)), []byte(reg.Name)); err != nil {
				return err
			}
		}
		return nil
	})
}

func FindInDB(db *bolt.DB, eui []byte) (regs []Registration, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return ErrNotInitialized
		}
		c := bucket.Cursor()
		for _, length := range []int{9, 7, 6} { // MA-S, MA-M and MA-L
			prefix := []byte(eui)[:length]
			for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(eui, k); k, v = c.Next() {
				regs = append(regs, Registration{Prefix: string(k), Name: string(v)})
			}
		}
		return nil
	})
	return regs, err
}

func InitializeDB(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(bucketName))
		_, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("Could not create bucket: %s", err)
		}
		return nil
	})
}
