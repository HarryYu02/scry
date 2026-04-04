package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

type BoltStore struct {
	db *bolt.DB
}

func openOrCreateBucket(tx *bolt.Tx, name string) (*bolt.Bucket, error) {
	bucket := tx.Bucket([]byte(name))
	if bucket == nil {
		var err error
		bucket, err = tx.CreateBucket([]byte(name))
		if err != nil {
			return nil, fmt.Errorf("create bucket: %s", err)
		}
	}
	return bucket, nil
}

func (b *BoltStore) GetIDTFMap(stem string) (map[string]float64, error) {
	var idTFMap map[string]float64
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("IDTF"))
		if bucket == nil {
			return fmt.Errorf("IDTF bucket not found")
		}

		idTFMapBytes := bucket.Get([]byte(stem))
		if idTFMapBytes == nil {
			return fmt.Errorf("idTFMapBytes not found")
		}
		err := json.Unmarshal(idTFMapBytes, &idTFMap)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return idTFMap, nil
}

func (b *BoltStore) GetTitle(docID string) (string, error) {
	var title string
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Title"))
		if bucket == nil {
			return fmt.Errorf("Title bucket not found")
		}

		titleBytes := bucket.Get([]byte(docID))
		if titleBytes == nil {
			return fmt.Errorf("title not found")
		}
		title = string(titleBytes)
		return nil
	})
	if err != nil {
		return "", err
	}
	return title, nil
}

func (b *BoltStore) GetWords() ([]string, error) {
	words := make([]string, 0)
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("WordCount"))
		if bucket == nil {
			return fmt.Errorf("WordCount bucket not found")
		}

		bucket.ForEach(func(k, v []byte) error {
			words = append(words, string(k))
			return nil
		})

		return nil
	})
	return words, nil
}

func (b *BoltStore) GetWordCount(word string) (int, error) {
	var count int
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("WordCount"))
		if bucket == nil {
			return fmt.Errorf("WordCount bucket not found")
		}

		countBytes := bucket.Get([]byte(word))
		if countBytes == nil {
			count = 0
			return nil
		}
		countNum, err := strconv.Atoi(string(countBytes))
		if err != nil {
			return err
		}
		count = countNum
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *BoltStore)GetPosition(docID string) (Position, error) {
	var pos Position
	err := b.db.View(func(tx *bolt.Tx) error {
		posBucket := tx.Bucket([]byte("Position"))
		if posBucket == nil {
			return fmt.Errorf("posBucket not found in index")
		}
		posBytes := posBucket.Get([]byte(docID))
		if posBytes == nil {
			return fmt.Errorf("cannot find %s in posBucket", docID)
		}

		err := json.Unmarshal(posBytes, &pos)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Position{}, err
	}
	return pos, nil
}

