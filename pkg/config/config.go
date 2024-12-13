package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	RaftAddress       string
	NodeID            string
	BootstrapNode     bool
	Port              string
	BadgerDBPath      string
	RaftDataDir       string
	Peers             []string
	SnapShotThreshold uint64
	SnapShotInterval  uint64
}

func LoadConfig() (*Config, error) {
	godotenv.Load(".env")

	raftAddr := os.Getenv("RAFT_ADDRESS")
	nodeID := os.Getenv("NODE_ID")
	peersStr := os.Getenv("RAFT_PEERS")
	bootstrapStr := os.Getenv("BOOTSTRAP")
	port := os.Getenv("PORT")
	dbPath := os.Getenv("BADGER_DB_PATH")
	snapshotThreshold := os.Getenv("SNAPSHOT_THRESHOLD")
	snapshotInterval := os.Getenv("SNAPSHOT_INTERVAL") // expected the value in seconds
	var snapshotThresholdInt, snapshotIntervalInt uint64
	if snapshotThreshold != "" {
		snapshotThresholdInt, _ = strconv.ParseUint(snapshotThreshold, 10, 32)
	} else {
		snapshotThresholdInt = 1000
	}

	if snapshotInterval != "" {
		snapshotIntervalInt, _ = strconv.ParseUint(snapshotInterval, 10, 32)
	} else {
		snapshotIntervalInt = 3600
	}

	var peers []string
	if peersStr != "" {
		peers = strings.Split(peersStr, ",")
	}
	return &Config{
		RaftAddress:       raftAddr,
		NodeID:            nodeID,
		BootstrapNode:     bootstrapStr == "true",
		Port:              port,
		BadgerDBPath:      dbPath,
		RaftDataDir:       "./raft-data",
		Peers:             peers,
		SnapShotInterval:  snapshotIntervalInt,
		SnapShotThreshold: snapshotThresholdInt,
	}, nil
}
