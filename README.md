# go-tcp-serv

A TCP server in Go implementing a simple SMTP-inspired command protocol with rate limiting and per-connection state management.

## Overview

This is a Go learning project exploring core networking and concurrency concepts:

- TCP connection handling with goroutines
- A stateful command protocol (clients must authenticate before accessing commands)
- IP-based rate limiting with mutex-protected state
- Read timeouts to prevent stalled connections

## Features

- Accepts multiple concurrent clients, each handled in its own goroutine
- Sends an SMTP-style welcome banner on connect (`220 <hostname>`)
- Stateful session: clients must send `EHLO` before accessing other commands
- IP-based rate limiting: max 5 concurrent connections per IP (returns `429 Too Many Requests` when exceeded)
- 10-second read timeout per message
- SMTP-style response codes (`250`, `221`, `500`, `550`)

## Protocol

| Command       | Requires EHLO | Description                        | Response                          |
|---------------|---------------|------------------------------------|-----------------------------------|
| `EHLO <name>` | No            | Handshake — sets session to live   | `250 Pleased to meet you <name>`  |
| `DATE`        | Yes           | Returns the current server time    | `250 <timestamp>`                 |
| `BYE`         | No            | Closes the connection gracefully   | `221 Bye`                         |
| *(unknown)*   | —             | Any unrecognised command           | `500 Unknown command`             |

## Getting Started

**Requirements:** Go 1.24+

```bash
# Clone and run
git clone <repo-url>
cd go-tcp-serv
go run ./cmd/main.go
```

The server listens on `localhost:8080`.

Connect with `nc` or `telnet`:

```bash
nc localhost 8080
```

## Example Session

```
< 220 localhost
> EHLO alice
< 250 Pleased to meet you alice
> DATE
< 250 2026-05-28 14:32:01.123456789 +0000 UTC
> BYE
< 221 Bye
```

If you try `DATE` before `EHLO`:

```
< 220 localhost
> DATE
< 550 Bad state
```

## Project Structure

```
go-tcp-serv/
├── cmd/
│   └── main.go                  # Entry point — starts the server on localhost:8080
├── internal/
│   ├── server/
│   │   └── server.go            # Accepts connections, dispatches commands
│   ├── client/
│   │   └── client.go            # Per-connection state (IsLive flag)
│   ├── config/
│   │   └── config.go            # Host/port configuration
│   ├── protocol/
│   │   ├── ehlo_command.go      # EHLO handler
│   │   ├── date_command.go      # DATE handler
│   │   └── bye_command.go       # BYE handler
│   └── ratelimiter/
│       └── rate_limiter.go      # IP-based connection limiter
└── go.mod
```
