package badger

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
)

var DB *badger.DB

type DBStorage interface {
	SaveEntries(id string, secrets []byte) error
	GetEntries(id string) (secrets []byte, err error)
	DeleteEntries(id string) error
	UpdateEntries(key string, value []byte) error
}

type Badger struct {
	raftNode *raft.Raft
}

func NewBadgerDBStorage(raftNode *raft.Raft) *Badger {
	return &Badger{raftNode: raftNode}
}

func (s *Badger) SaveEntries(key string, value []byte) error {
	cmd := Command{
		Op:    "store",
		Key:   key,
		Value: value,
	}

	data, err := cmd.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize command: %w", err)
	}

	applyFuture := s.raftNode.Apply(data, 5*time.Second)
	if err = applyFuture.Error(); err != nil {
		return fmt.Errorf("failed to apply Raft log: %w", err)
	}

	return nil
}

func (s *Badger) GetEntries(id string) ([]byte, error) {

	fullData, err := s.GetEntry(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file chunks: %v", err)
	}
	return fullData, nil
}

func (s *Badger) DeleteEntries(key string) error {
	cmd := Command{
		Op:  "delete",
		Key: key,
	}

	data, err := cmd.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize command: %w", err)
	}

	applyFuture := s.raftNode.Apply(data, 5*time.Second)
	if err = applyFuture.Error(); err != nil {
		return fmt.Errorf("failed to apply Raft log: %w", err)
	}

	return nil
}

func (s *Badger) UpdateEntries(key string, value []byte) error {
	cmd := Command{
		Op:    "update",
		Key:   key,
		Value: value,
	}

	data, err := cmd.Serialize()
	if err != nil {
		return fmt.Errorf("failed to serialize command: %w", err)
	}

	applyFuture := s.raftNode.Apply(data, 5*time.Second)
	if err = applyFuture.Error(); err != nil {
		return fmt.Errorf("failed to apply Raft log: %w", err)
	}

	return nil
}

func (s *Badger) GetEntry(key string) ([]byte, error) {
	var chunkData []byte

	err := DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return badger.ErrKeyNotFound
			}
			return fmt.Errorf("failed to get chunk for key %s: %v", key, err)
		}

		chunkData, err = item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("failed to retrieve chunk data for key %s: %v", key, err)
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			log.Printf("Chunk with key %s not found", key)
		} else {
			log.Printf("Error retrieving chunk with key %s: %v", key, err)
		}
		return nil, err
	}

	return chunkData, nil
}
