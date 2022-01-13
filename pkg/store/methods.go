package store

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cockroachdb/pebble"
)

func (s *Store) Start(addr string) {

	if s.cm != nil {
		go func() {
			for {
				entry := <-s.cm.CommitChan
				log.Printf("got committed entry %+v\n", entry)
				s.Set(entry.Key, entry.Value)
			}
		}()
	}

	http.HandleFunc("/get", s.handleGet())
	http.HandleFunc("/set", s.handleSet())
	log.Fatal(http.ListenAndServe(addr, nil))

}

//
// Get retrieves the value associated with a key from the database/map if it exists
// returns an error if it doesn't
//
func (s *Store) Get(key string) (string, error) {

	switch s.db {
	case nil:
		if value, ok := s.kv[key]; ok {
			return value, nil
		}
		return "", fmt.Errorf("key %s not in store", key)

	default:
		value, closer, err := s.db.Get([]byte(key))
		if err != nil {
			return "", err
		}
		closer.Close()
		return string(value), nil
	}
}

//
// Set adds the key/value pair to the map/ database unless a
// consensus module is present, in which case it sends the pair
// to the raft cluster to be replicated
//
func (s *Store) Set(key, value string) error {

	switch s.db {
	case nil:

		s.kv[key] = value
		log.Printf("set value in kv store %s : %s", key, value)
		return nil

	default:

		if err := s.db.Set([]byte(key), []byte(value), &pebble.WriteOptions{}); err != nil {
			return fmt.Errorf("error setting key %s to value %s -> %w", key, value, err)
		}
		log.Printf("set value in db %s : %s\n", key, value)
		return nil
	}
}
