

**Distributed File Storage in Go**

### Description
This project is a distributed file storage system built using Go. It supports file distribution across multiple nodes with fault tolerance and scalability. The system is designed to handle failures and ensures high availability by replicating files across nodes. It uses consistent hashing for efficient file storage and retrieval.

### Features
- Distributed architecture with file replication across nodes.
- Consistent hashing for balanced file distribution.
- Fault-tolerant storage with recovery mechanisms.
- Scalable design for adding/removing nodes dynamically.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/bhas10bc/distributed_file_storage_go.git
   ```
2. Navigate to the project directory:
   ```bash
   cd distributed_file_storage_go
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Configuration
Each node in the system can be configured using a `.env` file. You can customize the following settings:
- **NODE_PORT**: Define the port on which the node listens.
- **NODE_ID**: A unique identifier for each node in the distributed system.
- **REPLICATION_FACTOR**: The number of replicas for each file to ensure redundancy.

Example `.env` file:
```env
NODE_PORT=4000
NODE_ID=node1
REPLICATION_FACTOR=3
```

### Usage
1. Start each node by running:
   ```bash
   go run main.go
   ```
2. Use a client to interact with the nodes and upload/download files. The client sends a file to a node, which then distributes it according to the consistent hashing algorithm.

3. Example command to upload a file:
   ```bash
   curl -X POST -F 'file=@/path/to/file' http://localhost:4000/upload
   ```

### Architecture
The system uses consistent hashing to distribute files across nodes. Each file is hashed to determine its storage location. Files are replicated across multiple nodes based on the replication factor, ensuring that no single point of failure exists.

### License
This project is licensed under the MIT License.

### Contributors
- Bhaskar Chandran

---
