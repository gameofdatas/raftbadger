package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	storage "raftbadger/pkg/badger"
	"raftbadger/pkg/config"
	"raftbadger/pkg/factory"
	"raftbadger/pkg/proto"
	"raftbadger/pkg/raftinit"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
)

type Server struct {
	proto.BadgerServiceServer
	upSince       time.Time
	badgerStorage storage.DBStorage
	badgerStatus  bool
	RaftNode      *raft.Raft
	Tm            *transport.Manager
	fsm           *storage.FSM
	cfg           *config.Config
}

func NewServer(cfg *config.Config, isBootstrapNode bool) (*Server, error) {
	dbPath := os.Getenv("BADGER_DB_PATH")
	if dbPath == "" {
		dbPath = "./data"
	}

	err := factory.InitializeBadgerDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize BadgerDB: %v", err)
	}

	badgerStatus := storage.DB != nil
	raftNode, tm, fsm, err := factory.InitializeRaft(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Raft: %v", err)
	}
	if isBootstrapNode {
		log.Println("Bootstrapping the Raft cluster as this is the initial node.")
		err = raftinit.BootstrapSingleNode(raftNode, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to bootstrap Raft cluster: %v", err)
		}

		go func() {
			raftinit.WaitForLeaderElection(raftNode)
			raftinit.AddPeersIfLeader(raftNode, cfg)
		}()
	}

	return &Server{
		upSince:       time.Now(),
		badgerStorage: storage.NewBadgerDBStorage(raftNode),
		badgerStatus:  badgerStatus,
		RaftNode:      raftNode,
		Tm:            tm,
		fsm:           fsm,
		cfg:           cfg,
	}, nil
}

func (s *Server) SaveEntries(context context.Context, req *proto.SaveEntriesRequest) (*proto.SaveEntriesResponse, error) {
	id := req.GetId()
	secrets := req.GetSecrets()

	err := s.badgerStorage.SaveEntries(id, secrets)
	if err != nil {
		log.Printf("Failed to store data in BadgerDB: %v", err)
		return &proto.SaveEntriesResponse{
			Success: false,
		}, err
	}

	return &proto.SaveEntriesResponse{
		Success: true,
	}, nil
}

func (s *Server) GetEntries(context context.Context, req *proto.GetEntriesRequest) (*proto.GetEntriesResponse, error) {
	id := req.GetId()
	entries, err := s.badgerStorage.GetEntries(id)

	if err != nil {
		log.Printf("Failed to get data from BadgerDB: %v", err)
		return &proto.GetEntriesResponse{
			Id:   id,
			Data: nil,
		}, err
	}

	return &proto.GetEntriesResponse{
		Id:   id,
		Data: entries,
	}, nil
}

func (s *Server) DeleteEntries(context context.Context, req *proto.DeleteEntriesRequest) (*proto.DeleteEntriesResponse, error) {
	id := req.GetId()

	err := s.badgerStorage.DeleteEntries(id)

	if err != nil {
		log.Printf("Failed to get data from BadgerDB: %v", err)
		return &proto.DeleteEntriesResponse{
			Success: false,
		}, err
	}

	return &proto.DeleteEntriesResponse{
		Success: true,
	}, nil
}

func (s *Server) UpdateEntries(context context.Context, req *proto.UpdateEntriesRequest) (*proto.UpdateEntriesResponse, error) {
	id := req.GetId()
	secrets := req.GetSecrets()

	err := s.badgerStorage.UpdateEntries(id, secrets)
	if err != nil {
		log.Printf("Failed to update data in BadgerDB: %v", err)
		return &proto.UpdateEntriesResponse{
			Success: false,
		}, err
	}

	return &proto.UpdateEntriesResponse{
		Success: true,
	}, nil
}
