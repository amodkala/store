package store

import (
	"github.com/amodkala/raft"
	"github.com/cockroachdb/pebble"
)

type Store struct {
	kv map[string]string
	db *pebble.DB
	cm *raft.CM
}

func New(opts ...StoreOpt) *Store {
	s := &Store{}

	for _, opt := range opts {
		opt(s)
	}

	if s.db == nil {
		s.kv = make(map[string]string)
	}

	return s
}
