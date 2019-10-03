package db

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

// it's not regular init() which gets called automatically
func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func CreateTask(task string) (int, error) {
	// using this closure to have access to the id generated inside the func
	var id int

	err := db.Update(func(tx *bolt.Tx) error {
		// get the bucket
		b := tx.Bucket(taskBucket)

		// get the next id
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)

		// store in the bucket
		return b.Put(key, []byte(task))
	})

	if err != nil {
		return -1, err
	}

	return 0, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
