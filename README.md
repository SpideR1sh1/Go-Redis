# Build Redis From Scratch

This repository contains an implementation of a simplified Redis-like in-memory database built using Go. 
The project follows the step-by-step guide from "Build Redis from Scratch" and includes features such as 
RESP parsing, command execution, and data persistence. 

---

## **Features**

1. **In-Memory Key-Value Store**:
   - Implements basic Redis-like commands for managing strings and hashes.

2. **RESP Parsing**:
   - Built a custom parser for the Redis Serialization Protocol (RESP) to handle client requests and send responses.

3. **Concurrency**:
   - Utilized Go routines to support multiple simultaneous client connections efficiently.

4. **Data Persistence**:
   - Developed an Append Only File (AOF) mechanism to store commands and restore data upon server restart.

5. **Custom Server Implementation**:
   - Designed and implemented a custom server in Go to handle Redis-like operations.

---

## **Installation**

### **Prerequisites**
- Go (version 1.18 or later)
- Redis CLI (`redis-cli`)

### **Setup**
1. Clone the repository:
   ```bash
   git clone https://github.com/SpideR1sh1/redis-in-go.git
   cd redis-in-go
