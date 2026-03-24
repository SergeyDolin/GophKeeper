# GophKeeper

CLI password manager with client-server architecture.

## Features

- User registration & authentication
- Secure storage of secrets
- Client-side encryption (AES-GCM)
- CLI interface
- Multiple data types:
  - login/password
  - text
  - card
  - binary

## Architecture

Client encrypts data before sending to server.
Server stores only encrypted data.

## Run server

```bash
go run cmd/server/main.go