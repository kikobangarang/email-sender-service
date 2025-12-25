# email-sender-service

Lightweight email-sending service with an HTTP interface.
Supports asynchronous email sending via a worker pool, SQLite persistence, and SMTP delivery.

---

## Features

- **HTTP API** to queue emails (`POST /send/{email}`)
- **Worker pool** for asynchronous email sending
- **SQLite** storage for email jobs
- **SMTP** sender
- Supports **forms** and **json** content types
- Graceful shutdown with **signal handling**

---

## Getting Started

### 1. Clone the repository

```bash
git clone [https://github.com/kikobangarang/email-sender-service.git](https://github.com/kikobangarang/email-sender-service.git)
cd email-sender-service
go mod download
```

### 2. Set up environment variables

Create a `.env` file in the `cmd/server` directory with the following content(this is an example, replace with your actual SMTP server details):

```env
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your_username
SMTP_PASS=your_password
SMTP_FROM=your_email@example.com

```

### 3. Run the server

```bash
cd cmd/server
go run main.go
```

### 4. Check the server is running

You can use `curl` to check if the server is running:

```bash
curl http://localhost:8080/health
```

A 200 OK response indicates the server is running.

### 5. Send a test email

You can use `curl` to send a test email:

```bash
curl -X POST http://localhost:8080/send/{destination_email} \
     -H "Content-Type: application/json" \
     -d '{"subject":"Test Email","body":"This is a test email."}'
```

A 202 Accepted response indicates the email has been queued for sending.

---

### Configuration

The worker pool size and other settings can be adjusted in the `main.go` file, using the following struct.

```go
	workerCfg := email.WorkerConfig{
		WorkerCount:  3,
		PollInterval: 2 * time.Second,
		MaxRetries:   3,
		BatchSize:    10,
	}
```
