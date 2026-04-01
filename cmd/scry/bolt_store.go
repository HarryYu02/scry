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

func (b *BoltStore) GetIDTFMap(stem string) (map[string]float64, error) {
	var idTFMap map[string]float64
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("IDTF"))
		if bucket == nil {
			return fmt.Errorf("IDTF bucket not found")
		}

		idTFMapBytes := bucket.Get([]byte(stem))
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

		bucket.ForEach(func (k, v []byte) error {
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
