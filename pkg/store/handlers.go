package store

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/amodkala/raft"
)

type getResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *Store) handleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		keys, ok := r.URL.Query()["key"]
		if !ok {
			http.Error(w, "no value for \"key\" parameter", http.StatusBadRequest)
			return
		}

		results := []getResponse{}
		for _, key := range keys {
			if value, err := s.Get(key); err != nil {
				result := getResponse{
					Key:   key,
					Value: value,
				}
				results = append(results, result)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		jsonResults, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Write(jsonResults)

	}
}

type setRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *Store) handleSet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var reqBody setRequest

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		key, value := reqBody.Key, reqBody.Value

		log.Printf("got kv pair %s : %s\n", key, value)

		switch s.cm {
		case nil:
			if err := s.Set(key, value); err != nil {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "successfully set kv pair %s : %s", key, value)
				return
			}
		default:
			entry := raft.Entry{
				Key:   key,
				Value: value,
			}
			if replicated := s.cm.Replicate([]raft.Entry{entry}); replicated {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "%s : %s replicated\n", key, value)
				return
			}
			http.Error(w, fmt.Sprintf("%s : %s couldn't be replicated\n", key, value), http.StatusBadRequest)
			return
		}
	}
}
