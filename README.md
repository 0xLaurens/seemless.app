# BitDash - bitda.sh 

BitDash is a Airdrop like application which allows users to share files with eachother no matter the platform the end user uses. The application works peer-to-peer to give the users exchanging files privacy, the server doesn't know what was exchanged since it never stored any of the files.

## Features
- [ ] p2p connections
- [ ] user discovery 
- [ ] local network file sharing
- [ ] file broadcasting to all users
- [ ] outside of local-network file sharing
- [ ] accept/deny files received 

## Backend
The backend application is a Go application that serves as a signaling mechanism for the webRTC process. The server facilitates communication between two or more devices in a webRTC session. The server acts as the middleman exchanging the information that is necessary for establisihing and managing peer-to-peer connections. The server helps devices discover each other, negotiate session details, and exchange signaling messages such as offer, answer and ICE (Interactive Connectivity Establishment) candidates.

### Installation
***Backend*** requirements:
- [golang 1.20](https://go.dev/doc/install)  or newer

```sh
git clone https://gitlab.com/ihomer/ihomer-academy/stage/laurensdrop
go install
```


### Building
```sh
go build
```
### Hot-Reloading 
```
air
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
The user interface for the application. Written in Vite + Vue using typescript. Helps clients share files. The user interface is a PWA allowing users to install it like a app. 

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
