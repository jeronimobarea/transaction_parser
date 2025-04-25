# Transaction Parser

An Ethereum transaction parser that enables querying and monitoring of incoming and outgoing transactions for subscribed addresses. This tool provides a simple interface (CLI or HTTP API) to:

- Retrieve the latest block number
- Subscribe an address for transaction monitoring
- Fetch inbound and outbound transactions for any subscribed address

---

## ⚙️ Features

- **Blockchain Interaction**: Uses Ethereum JSON-RPC to interact with any EVM-compatible node.
- **Lightweight Storage**: In-memory storage by default, easily swappable for persistent backends.
- **Pure Go**: No external dependencies beyond the standard library.
- **Modular Design**: Clear separation of parser, repository, client, and HTTP handlers.

---

## 🚀 Getting Started

### Prerequisites
- Go 1.24+
- Docker (optional)

### Installation

```bash
# Clone the repo
git clone https://github.com/jeronimobarea/transaction_parser.git
cd transaction_parser
```

### Running via Docker

```bash
make run
```

### Running Locally

```bash
# Start the server
go run cmd/main.go
```

### Run tests
```bash
make tests
```

# By default, the service listens on port 3000

---

## 🔌 API Reference

All endpoints respond with `Content-Type: application/json`.

### 1. Get Current Block

```curl
curl --location 'http://localhost:3000/blocks/current'
```

#### Response

```json
{
  "currentBlock": 22347822
}
```

### 2. Subscribe Address

```
curl --location 'http://localhost:3000/subscribe?address=<YOUR_ADDRESS>'
```

- **[REQUIRED] Query Parameter**: `address` — EVM-compatible address (0x-prefixed, 40 hex chars)

#### Response
- **Status**: `200 OK` on success

### 3. Get Transactions

```
curl --location 'http://localhost:3000/transactions?address=<YOUR_ADDRESS>'
```
- **[REQUIRED] Query Parameter**: `address` — EVM-compatible address (0x-prefixed, 40 hex chars)

#### Response
```json
[
  {
    "hash": "0x5140...e020",
    "from": "0x4838...d9ee7",
    "to":   "0xe688...7127",
    "value":"0x2bf5fe4aff5181",
    "blockNumber":"0x1550035"
  }
]
```

---

## 🗂️ Project Structure

```
├── cmd/                              # Entry points
│   └── main.go                       # Main application
├── internal/                         # Internal code
│   ├── platform/                     # API setup
│   ├── chains/ethereum/              # Ethereum-specific parserr, poller & client
│   ├── chains/ethereum/repository/   # In-memory storage implementation 
│   ├── parser/                       # Parser interface, repository, handlers
│   ├── test/                         # Mock implementations and helpers
│   ├── pkg/svcerrors/                # Shared common service errors
│   ├── pkg/httphandler/              # HTTP Handler util
│   └── pkg/evm/                      # EVM address utils
├── pkg/osx                           # Shared libraries OSX
```

---

## 📝 Notes

- **Complexity**: I've added some additional complexity like _sync.Map_ in order to show an example of modularity for the excercise, for a real case scenario adding that type of modular complexity too early can be counterproductive.
- **Testing**: Currently uses unit tests with handwritten fakes. In a real word case, I would use integration tests and mocking frameworks (e.g., `mockgen`, `testify`, `mmock`), and added extended coverage for edge cases.
- **Improvements**:
  - EIP-55 checksum validation for addresses
  - Rate limiting / batching of RPC calls
  - Graceful shutdown and health checks
  - Proper route http method setting, and router robustness
  - Improve error handling in the service and handlers
  - Better logging with more modulatiry
  And more.
---

