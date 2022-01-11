package store

import (
	"fmt"
	"log"

	"github.com/cockroachdb/pebble"
)

func (s *Store) Get(key []byte) ([]byte, error) {

	switch s.db {
	case nil:
		return []byte(s.kv[string(key)]), nil
	default:
		value, closer, err := s.db.Get(key)
		defer closer.Close()
		if err != nil {
			return nil, err
		}
		return value, err
	}
}

func (s *Store) Set(key, value []byte) error {

	switch s.db {
	case nil:
		s.kv[string(key)] = string(value)
		log.Printf("set value in kv store %s : %s", key, value)
		return nil
	default:
		/*
			do consensus client interaction here
		*/
		if err := s.db.Set(key, value, &pebble.WriteOptions{}); err != nil {
			return fmt.Errorf("error setting key %s to value %s -> %w", key, value, err)
		}
		log.Printf("set value in db %s : %s", key, value)
		return nil
	}
}
