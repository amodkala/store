package store

import (
	"log"

	"github.com/amodkala/raft"
	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
)

type StoreOpt func(*Store)

func WithCM(cm *raft.CM) func(*Store) {
	return func(s *Store) {
		s.cm = cm
	}
}

func WithDB() func(*Store) {
	return func(s *Store) {
		db, err := pebble.Open("", &pebble.Options{FS: vfs.NewMem()})
		if err != nil {
			log.Printf("couldn't add db to store -> %v\n", err)
			return
		}
		s.db = db
	}
}
