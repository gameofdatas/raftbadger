package raftinit

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	storage "raftbadger/pkg/badger"
	"raftbadger/pkg/config"
	"raftbadger/pkg/utils"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
)

func InitializeRaft(cfg *config.Config) (*raft.Raft, *transport.Manager, *storage.FSM, error) {
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(cfg.NodeID)
	raftConfig.CommitTimeout = 1 * time.Second

	raftConfig.SnapshotThreshold = cfg.SnapShotThreshold
	raftConfig.SnapshotInterval = time.Duration(cfg.SnapShotInterval) * time.Second
	raftConfig.MaxAppendEntries = 1000

	baseDir := filepath.Join(cfg.RaftDataDir, cfg.NodeID)
	if err := utils.EnsureDir(baseDir); err != nil {
		return nil, nil, nil, fmt.Errorf("node %s failed to ensure directory: %v", cfg.NodeID, err)
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("node %s failed to create log store: %v", cfg.NodeID, err)
	}

	sdb, err := raftboltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf(`boltdb.NewBoltStore(%q): %v`, filepath.Join(baseDir, "stable.dat"), err)
	}
	snapshotStore, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("node %s failed to create snapshot store: %v", cfg.NodeID, err)
	}

	tm := transport.New(raft.ServerAddress(cfg.RaftAddress), []grpc.DialOption{grpc.WithInsecure()})

	fsm := &storage.FSM{}

	raftNode, err := raft.NewRaft(raftConfig, fsm, logStore, sdb, snapshotStore, tm.Transport())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("node %s failed to initialize Raft: %v", cfg.NodeID, err)
	}

	return raftNode, tm, fsm, nil
}

func BootstrapSingleNode(raftNode *raft.Raft, config *config.Config) error {
	cfg := raft.Configuration{
		Servers: []raft.Server{
			{
				Suffrage: raft.Voter,
				ID:       raft.ServerID(config.NodeID),
				Address:  raft.ServerAddress(config.RaftAddress),
			},
		},
	}

	// Bootstrap this node as the leader
	future := raftNode.BootstrapCluster(cfg)
	if err := future.Error(); err != nil {
		return fmt.Errorf("failed to bootstrap Raft cluster: %v", err)
	}
	log.Printf("Raft node %s bootstrapped as a single-node cluster", config.NodeID)
	return nil
}

func AddPeersIfLeader(raftNode *raft.Raft, cfg *config.Config) {
	if raftNode.State() != raft.Leader {
		log.Println("This node is not the leader. Skipping peer addition.")
		return
	}

	for _, peer := range cfg.Peers {
		peerParts := strings.Split(peer, ":")
		if len(peerParts) != 2 {
			log.Printf("Invalid peer address: %s", peer)
			continue
		}

		peerID := peerParts[0]
		port := peerParts[1]
		peerAddress := fmt.Sprintf("localhost:%s", port)
		// Attempt to add the peer with retry logic
		err := addPeerIfLeader(raftNode, peerID, peerAddress, 5)
		if err != nil {
			log.Printf("Error adding peer %s: %v", peer, err)
		} else {
			log.Printf("Successfully added peer %s", peerID)
		}
	}
}

func WaitForLeaderElection(raftNode *raft.Raft) {
	for {
		if raftNode.State() == raft.Leader {
			log.Println("Node has been elected leader. Proceeding to add peers.")
			break
		}
		log.Println("Waiting for node to be elected leader...")
		time.Sleep(2 * time.Second)
	}
}

func addPeerIfLeader(raftNode *raft.Raft, peerID, peerAddress string, maxRetries int) error {
	var attempt int

	for attempt = 0; attempt < maxRetries; attempt++ {

		// Try to add a voter (peer) to the Raft cluster
		future := raftNode.AddVoter(raft.ServerID(peerID), raft.ServerAddress(peerAddress), 0, 0)
		if future.Error() == nil {
			log.Printf("Successfully added peer %s at %s", peerID, peerAddress)
			return nil // Success
		}

		// Log the failure and retry after a backoff
		log.Printf("Failed to add peer %s: %v (attempt %d)", peerID, future.Error(), attempt+1)

		// Exponential backoff: wait 2^attempt seconds (with jitter)
		backoffDuration := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		time.Sleep(backoffDuration)
	}

	return fmt.Errorf("failed to add peer after %d attempts", maxRetries)
}
