package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"raftbadger/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Connect to the gRPC server
func connectToNode(address string) (*grpc.ClientConn, proto.BadgerServiceClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to node at %s: %v", address, err)
	}
	client := proto.NewBadgerServiceClient(conn)
	return conn, client, nil
}

// Write data to the leader
func writeToLeader(leaderAddress, id string, data []byte) error {
	conn, client, err := connectToNode(leaderAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &proto.SaveEntriesRequest{
		Id:      id,
		Secrets: data,
	}
	res, err := client.SaveEntries(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to save data: %v", err)
	}
	if res.Success {
		log.Println("Data successfully written to the leader.")
	} else {
		log.Println("Failed to write data to the leader.")
	}
	return nil
}

// Read data from a follower
func readFromFollower(followerAddress, id string) ([]byte, error) {
	conn, client, err := connectToNode(followerAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &proto.GetEntriesRequest{
		Id: id,
	}
	res, err := client.GetEntries(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %v", err)
	}

	log.Printf("Data retrieved from follower: %s\n", string(res.Data))
	return res.Data, nil
}

func main() {

	err := writeToLeader("localhost:50051", "test_key1", []byte("hello from leader14"))
	if err != nil {
		log.Fatalf("Error writing to leader: %v", err)
	}
	time.Sleep(2 * time.Second)
	followerAddresses := []string{"localhost:50052", "localhost:50053"}
	for _, followerAddr := range followerAddresses {
		_, err = readFromFollower(followerAddr, "test_key1")
		if err != nil {
			log.Printf("Error reading from follower %s: %v", followerAddr, err)
		}
	}

}
