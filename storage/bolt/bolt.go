package bolt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slackbot/storage"

	"github.com/boltdb/bolt"
)

const bucket = "slack"

type Bolt struct {
	db   *bolt.DB
	Path string
}

func (b *Bolt) Add(user string, msg string) {
	err := errors.New("tmp")
	b.db, err = bolt.Open(b.Path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer b.db.Close()

	var v []byte
	var entries []string
	b.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(bucket))
		if buck == nil {
			return nil
		}
		v = buck.Get([]byte(user))
		return nil
	})
	if string(v) != "" {
		if err := json.Unmarshal(v, &entries); err != nil {
			panic(err)
		}
	}

	encoded, err := json.Marshal(append(entries, msg))
	if err != nil {
		panic(err)
	}

	updateErr := b.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(user), []byte(encoded))
	})

	if updateErr != nil {
		fmt.Printf("Error putting %s in key %s; %s\n", msg, user, updateErr)
		panic(updateErr)
	}
}

func (b *Bolt) Get(user string) string {
	err := errors.New("tmp")
	b.db, err = bolt.Open(b.Path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer b.db.Close()

	var v []byte
	var entries []string
	b.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(bucket))
		if buck == nil {
			return nil
		}
		v = buck.Get([]byte(user))
		return nil
	})
	if string(v) != "" {
		if err := json.Unmarshal(v, &entries); err != nil {
			panic(err)
		}
		return entries[storage.Random(len(entries))]
	}
	return storage.Nothing
}
