package factory

import (
	"fmt"
	"os"
	storage "raftbadger/pkg/badger"
	"raftbadger/pkg/config"
	"raftbadger/pkg/raftinit"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/dgraph-io/badger/v4"
	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
)

// InitializeBadgerDB encapsulates BadgerDB initialization.
func InitializeBadgerDB(dbPath string) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dbPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create BadgerDB directory: %w", err)
		}
	}

	opts := badger.DefaultOptions(dbPath).
		WithLoggingLevel(badger.INFO). // Set logging level for BadgerDB
		WithSyncWrites(true)           // Enable sync writes for durability

	db, err := badger.Open(opts)
	if err != nil {
		return fmt.Errorf("failed to open BadgerDB at path %s: %w", dbPath, err)
	}

	storage.DB = db
	return nil
}

// InitializeRaft encapsulates Raft initialization logic.
func InitializeRaft(config *config.Config) (*raft.Raft, *transport.Manager, *storage.FSM, error) {
	return raftinit.InitializeRaft(config)
}

// InitializeGRPCServer encapsulates the creation of a gRPC server.
func InitializeGRPCServer() *grpc.Server {
	return grpc.NewServer()
}
