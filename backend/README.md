# Number Finder API

A high-performance REST service that efficiently finds index positions in a sorted number sequence using binary search.

---

## 🚀 Features

- Fast binary search implementation for large sorted datasets
- RESTful API using Fiber framework
- Environment-based configuration
- Structured logging
- Graceful shutdown

---

## 🛠 Prerequisites

- Go 1.23 or higher
- Make (optional, for using Makefile commands)

---

## 🔧 Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/igor-silveira/number-finder.git
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

4. Modify the `.env` file according to your needs:
   ```env
   PORT=8080
   LOG_LEVEL=info
   DATA_PATH=data/input.txt
   ```

---

## ▶️ Usage

1. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

2. The API will be available at `http://localhost:8080` (or your configured port).

---

## 📖 API Endpoints

### **GET** `/api/number/:value`

Find the index of a number in the sorted sequence.

#### Path Parameters

- `value` (number, required): The number to find in the sequence

#### Query Parameters

- `thresholdPercentage` (float, optional): Percentage threshold for approximate matching
  - If provided, the service will try to find a value within this percentage threshold
  - Default value: `0.0`

**Example Request:**

```http
GET /api/number/42?thresholdPercentage=0.1
```

**Success Response:**

- **Code:** 200 OK
- **Content:**
  ```json
  {
      "index": 3,
      "value": 100,
      "is_approximate": false
  }
  ```

**Error Responses:**

- **Code:** 400 Bad Request
  ```json
  {
      "error": "Invalid value parameter"
  }
  ```

- **Code:** 404 Not Found
  ```json
  {
      "error": "No value found within acceptable threshold"
  }
  ```

**Example Command:**
```bash
curl http://localhost:8080/api/number/100
```

---

## ⚙️ Configuration

The service can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level [debug, info, error] (default: info)
- `DATA_PATH`: Path to the data file containing the sorted sequence

---

## 🧪 Running Tests

Run the test suite using Make:
```bash
make test
```

---

## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.