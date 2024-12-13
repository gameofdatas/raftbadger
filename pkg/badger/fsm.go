package badger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"raftbadger/pkg/config"

	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
)

type FSM struct {
	mu   sync.Mutex
	conf config.Config
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	panic("implement me")
}

func (f *FSM) Restore(snapshot io.ReadCloser) error {
	panic("implement me")
}

func (f *FSM) Apply(log *raft.Log) interface{} {
	var cmd Command

	if err := cmd.Deserialize(log.Data); err != nil {
		return nil
	}

	switch cmd.Op {
	case "store":
		return f.storeEntry(cmd.Key, cmd.Value)
	case "update":
		return f.updateEntry(cmd.Key, cmd.Value)
	case "delete":
		return f.deleteEntry(cmd.Key)
	default:
		return nil
	}
}

func (f *FSM) storeEntry(key string, value []byte) error {
	log.Printf("Attempting to store entry with key: %s", value)

	// Store the serialized data in BadgerDB
	err := DB.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(key), value); err != nil {
			log.Printf("BadgerDB transaction error for entry with key: %s. Error: %v", key, err)
			return err
		}
		log.Printf("Successfully stored entry with key: %s", key)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to store chunk: %v", err)
	}

	return nil
}

func (f *FSM) updateEntry(id string, value []byte) error {
	err := DB.Update(func(txn *badger.Txn) error {

		_, err := txn.Get([]byte(id))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return fmt.Errorf("entry with ID %s not found", id)
			}
			return fmt.Errorf("failed to get entry: %w", err)
		}

		err = txn.Set([]byte(id), value)
		if err != nil {
			return fmt.Errorf("failed to update entry: %w", err)
		}

		return nil
	})

	return err
}

func (f *FSM) deleteEntry(key string) error {
	return DB.Update(func(txn *badger.Txn) error {

		_, err := txn.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("failed to get value with key %s: %w", key, err)
		}

		err = txn.Delete([]byte(key))
		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			return fmt.Errorf("failed to delete base entry: %w", err)
		}

		err = txn.Delete([]byte(key))
		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			return fmt.Errorf("failed to delete base entry: %w", err)
		}

		return nil
	})
}
