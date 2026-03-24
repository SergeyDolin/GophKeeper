# GophKeeper

GophKeeper is a CLI-based password manager with a client-server architecture written in Go.

The system allows users to securely store and manage sensitive data such as credentials, text notes, binary files, and bank card information.

## ✨ Features

* User registration and authentication (JWT)
* Secure password hashing (bcrypt)
* Client-side encryption (AES-256-GCM)
* Zero-knowledge architecture (server cannot read user data)
* Cross-platform CLI (Windows, Linux, macOS)
* Multiple data types:

  * Login / Password
  * Text notes
  * Binary data (files)
  * Bank card data
* Metadata support for all entries
* Unit tests (core logic covered)

---

## Architecture

```
CLI Client
   ↓
Encrypts data (AES-GCM)
   ↓
HTTP API
   ↓
Server stores encrypted data only
```

### Security model

* Passwords are hashed using **bcrypt**
* Encryption key is derived using **Argon2**
* Data is encrypted on the client before sending
* Server stores only encrypted data

Even if the database is compromised, attacker cannot read user data.

---

## Getting Started

### 1. Clone repository

```bash
git clone https://github.com/SergeyDolin/GophKeeper.git
cd GophKeeper
```

---

### 2. Run server

```bash
go run cmd/server/main.go
```

Server will start on:

```
http://localhost:8080
```

---

### 3. Run client

```bash
go run cmd/client/main.go login
```

---

## CLI Usage

### Login

```bash
go run cmd/client/main.go login
```

---

### Add secret

```bash
go run cmd/client/main.go add
```

You will be prompted for:

* Type (text / login / card / binary)
* Data
* Metadata
* Master password

---

### List secrets

```bash
go run cmd/client/main.go list
```

---

### Show version

```bash
go run cmd/client/main.go version
```

---

## Supported Data Types

| Type   | Description           |
| ------ | --------------------- |
| login  | Login + password pair |
| text   | Plain text notes      |
| binary | File content          |
| card   | Bank card data        |

---

## Testing

Run all tests:

```bash
go test ./...
```

Check coverage:

```bash
go test ./... -cover
```

---

## Documentation

Generate documentation:

```bash
go doc ./...
```

Run web documentation:

```bash
godoc -http=:6060
```

Open in browser:

```
http://localhost:6060
```

---

## Project Structure

```
cmd/
  server/        # server entrypoint
  client/        # CLI entrypoint

internal/
  handlers/      # HTTP handlers
  service/       # business logic
  storage/       # in-memory storage
  models/        # data models
  clientapi/     # HTTP client
  clientcrypto/  # encryption logic
  cli/           # CLI commands
```

---

## 👨‍💻 Author

Sergei Dolin

---