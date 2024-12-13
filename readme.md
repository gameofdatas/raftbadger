# RaftBadger Service

The **RaftBadger Service** is a distributed key-value store that combines Raft consensus for replication and BadgerDB for efficient storage. This project demonstrates how to implement a fault-tolerant storage system using gRPC for communication and Raft for cluster management.

---

## **Project Structure**

The project is organized as follows:

```plaintext
.
├── go.mod                       # Go module file
├── go.sum                       # Dependency lock file
├── main.go                      # Entry point of the application
├── pkg/
│   ├── badger/                  # BadgerDB-related logic
│   │   ├── badger.go            # Core BadgerDB operations
│   │   ├── command.go           # Command management for Raft FSM
│   │   └── fsm.go               # Finite State Machine implementation for Raft
│   ├── config/
│   │   └── config.go            # Configuration loader and management
│   ├── factory/
│   │   └── factory.go           # gRPC server initialization logic
│   ├── proto/                   # Protocol Buffers definitions and generated code
│   │   ├── badger.pb.go         # Generated Go code for gRPC services
│   │   └── badger.proto         # Protocol Buffers definition file
│   └── raftinit/
│       └── raftinit.go          # Raft initialization and setup
└── server/
└── server.go                # gRPC server implementation and service handlers
```

---

## **Features**

- **Raft Consensus**: Ensures fault-tolerant replication of data across nodes.
- **BadgerDB Integration**: Provides a fast and efficient key-value storage engine.
- **gRPC API**: Enables seamless communication for saving, retrieving, updating, and deleting data.
- **Scalability**: Designed to support multiple nodes with leader election and state replication.
- **Extensible Architecture**: Modular codebase for easy integration of new features.

---

## **Getting Started**

### **Prerequisites**
- Go 1.18+ installed.
- Protocol Buffers compiler (`protoc`) installed. [Installation Guide](https://grpc.io/docs/protoc-installation/)
- Git installed.

### **Setup**
1. Clone the repository:
```bash
git clone <repository-url>
    cd raftbadger
```
2. Install dependencies:
```bash
go mod tidy
```

3. Generate Protocol Buffers code (if not already generated):
```bash
protoc --proto_path=pkg/proto --go_out=pkg/proto --go-grpc_out=pkg/proto pkg/proto/badger.proto
```

4. Build the application:
```bash
go build -o raftbadger .
```

5. Run the application:
```bash
# run the followers replica
> terminal 1
NODE_ID=replica2 RAFT_ADDRESS=localhost:50052 BADGER_DB_PATH=./data/replica2 PORT=50052 BOOTSTRAP=false go run main.go
# run the second follower consider its a 3 replica raft cluster
> terminal 2
NODE_ID=replica3 RAFT_ADDRESS=localhost:50053 BADGER_DB_PATH=./data/replica3 PORT=50053 BOOTSTRAP=false go run main.go

# run the leader node with bootstrap and follower address as peers
>terminal 3
 NODE_ID=replica1 RAFT_ADDRESS=localhost:50051 RAFT_PEERS=replica2:50052,replica3:50053 BADGER_DB_PATH=./data/replica1 PORT=50051 BOOTSTRAP=true go run main.go
 
```

6. Run the client via [client](client/client.go)

```go

# the code writes to leader and reads from followers considering the RAFT replicates the record to followers successfully

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
```

---

## **Configuration**

The configuration is managed via the `pkg/config/config.go` file. You can customize settings such as:
- **Port**: The gRPC server port.
- **Bootstrap Node**: Whether the node is a Raft bootstrap node.

Example configuration file (`config.json`):
```json
{
"Port": "50051",
"BootstrapNode": true,
"RaftDirectory": "/tmp/raft",
"BadgerDirectory": "/tmp/badger"
}
```

---

## **API Overview**

The gRPC API exposes the following services:

### **Service: BadgerService**

1. **SaveEntries**
- Saves a key-value pair in the system.
- Request: `SaveEntriesRequest { id: string, secrets: bytes }`
- Response: `SaveEntriesResponse { success: bool }`

2. **GetEntries**
- Retrieves a value by its key.
- Request: `GetEntriesRequest { id: string }`
- Response: `GetEntriesResponse { id: string, data: bytes }`

3. **DeleteEntries**
- Deletes a key-value pair.
- Request: `DeleteEntriesRequest { id: string }`
- Response: `DeleteEntriesResponse { success: bool }`

4. **UpdateEntries**
- Updates an existing key-value pair.
- Request: `UpdateEntriesRequest { id: string, secrets: bytes }`
- Response: `UpdateEntriesResponse { success: bool }`
---

## **Development**

### **Project Modules**
- **BadgerDB Logic (`pkg/badger/`)**:
- Handles interactions with BadgerDB and data persistence.
- **Raft Initialization (`pkg/raftinit/`)**:
- Sets up Raft nodes, handles leader election, and state replication.
- **Protocol Buffers (`pkg/proto/`)**:
- Defines gRPC services and message formats.
- **Server Implementation (`server/`)**:
- Implements the gRPC services and registers them with the gRPC server.

---


## **Contributing**

We welcome contributions! Please follow these steps:
1. Fork the repository.
2. Create a feature branch:
```bash
git checkout -b feature/your-feature
```
3. Commit changes:
```bash
git commit -m "Add your feature"
```
4. Push and create a pull request.

---
