# Seemless.app

Seemless.app is a peer-to-peer file sharing application. It allows users to share files with each other without the need
of a central server. The application is a PWA allowing users to install it like a native application.

## Features

- [x] p2p connections
- [x] user discovery
- [x] local network file sharing
- [x] send files to a single user
- [x] send files to multiple users
- [x] accept/deny files received
- [x] outside of local-network file sharing

## Backend

The backend application is a Go application that serves as a signaling mechanism for the webRTC process. The server
facilitates communication between two or more devices in a webRTC session. The server acts as the middleman exchanging
the information that is necessary for establisihing and managing peer-to-peer connections. The server helps devices
discover each other, negotiate session details, and exchange signaling messages such as offer, answer and ICE (
Interactive Connectivity Establishment) candidates.

### Installation

***Backend*** requirements:

- [golang 1.20](https://go.dev/doc/install)  or newer

```sh
git clone https://gitlab.com/ihomer/ihomer-academy/stage/laurensdrop
```

### Running

```sh
go run cmd/web/main.go #replace cmd/web/main.go with other adapter if they exist
```

### Building

```sh
go build cmd/web/main.go
```

### Hot-Reloading

```sh
air --build.cmd "go build -o ./bin/api ./cmd/web/main.go" --build.bin "./bin/api"

```

### Test

```sh 
go test ./...
```

### Test Coverage

```sh
go test -coverprofile=coverage.out ./... ;  go tool cover -html=coverage.out
```

## Frontend

The user interface for the application. Written in Vite + Vue using typescript. Helps clients share files. The user
interface is a PWA allowing users to install it like a app.

### Installation

***Frontend*** requirements:

- [node](https://nodejs.org/en/download)

```sh
npm i
```

### Building

```sh
npm run build
```

### Hot-Reloading

```sh 
npm run build
```

### Unit testing

```
npm run test:unit 
```

### e2e testing

```sh 
npm run test:e2e
```
