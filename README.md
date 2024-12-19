
# Build Redis From Scratch

This repository contains an implementation of a simplified Redis-like in-memory database built using Go. 
The project follows a step-by-step guide and covers building a server, RESP parsing, implementing commands, and adding data persistence.

---

## **Features**

1. **Custom Server**:
   - Implements a custom TCP server to handle Redis-like client connections and requests.

2. **RESP Parsing**:
   - Built both a RESP reader and writer to handle Redis Serialization Protocol for processing client commands and sending responses.

3. **Command Execution**:
   - Supports core Redis commands for strings (`SET`, `GET`) and hashes (`HSET`, `HGET`).

4. **Concurrency**:
   - Utilized Go routines to handle multiple simultaneous client connections.

5. **Data Persistence**:
   - Developed an Append Only File (AOF) mechanism to log executed commands and restore the state upon server restart.

6. **Scalable Design**:
   - Modular code structure with individual files for server logic, RESP parsing, and AOF handling.

---

## **Installation**

### **Prerequisites**
- Go (version 1.18 or later)
- Redis CLI (`redis-cli`)

### **Setup**
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/build-redis-from-scratch.git
   cd build-redis-from-scratch
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. Interact with the server using `redis-cli`:
   ```bash
   redis-cli -h 127.0.0.1 -p 6379
   ```

---

## **Implemented Commands**

- **String Commands**:
  - `SET key value`: Store a key-value pair.
  - `GET key`: Retrieve the value associated with a key.

- **Hash Commands**:
  - `HSET key field value`: Set a field in a hash.
  - `HGET key field`: Get the value of a field in a hash.

---

## **Project Structure**

```
.
├── aof.go         # Handles data persistence with Append Only File (AOF)
├── handler.go     # Contains logic to process client commands
├── main.go        # Entry point for the Redis server
├── resp.go        # Implements RESP parsing (reader and writer)
└── README.md      # Documentation
```

---

## **How It Works**

### 1. **RESP Parsing**
- The server uses a custom RESP reader (`resp.go`) to decode incoming commands and a writer to encode responses.

### 2. **Concurrency**
- Uses Go routines to manage multiple client connections efficiently.

### 3. **Data Persistence**
- Implements AOF to store all executed commands for recovery during a server restart.

---

## **Usage Examples**

### Setting and Getting a Key
```bash
redis-cli -h 127.0.0.1 -p 6379
127.0.0.1:6379> SET name "RedisClone"
OK
127.0.0.1:6379> GET name
"RedisClone"
```

### Using Hashes
```bash
redis-cli -h 127.0.0.1 -p 6379
127.0.0.1:6379> HSET user name "Alice"
1
127.0.0.1:6379> HGET user name
"Alice"
```

---

## **Planned Features**

- Add support for additional Redis data types (e.g., lists, sets, sorted sets).
- Improve concurrency with better locking mechanisms.
- Extend RESP parsing for pipelined commands.

---

## **License**

This project is licensed under the MIT License.
