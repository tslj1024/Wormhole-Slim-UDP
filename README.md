# Wormhole-Slim-UDP

![Static Badge](https://img.shields.io/badge/license-GPLv3.0-blue)

**English** | [中文](README.zh_CN.md)

A lightweight UDP intranet penetration proxy

## Introduction

Wormhole-Slim-UDP is a lightweight UDP intranet penetration proxy tool designed to help users easily expose and access intranet services.

To use this software, you need to understand the following:

1. "Wormhole" translates to "虫洞" in Chinese. This name stems from the similarity between intranet penetration and wormholes in space.
2. This software is a lightweight version, primarily providing traffic forwarding.
3. Your intranet services are based on the UDP protocol.

## Features

- **Lightweight Design**: Streamlined core functionality and low resource usage
- **UDP Traversal**: Supports traversal of any UDP-based service
- **Simple Configuration**: Easy to configure and use
- **Cross-Platform**: Developed in the Go language, supports mainstream operating systems such as Windows, Linux, and macOS
- **Multi-Client**: Uses ClientIDs to distinguish multiple clients
- **Client Reconnect**: When the server crashes and restarts, the client can automatically reconnect

## Installation Method

### Binary Installation

Download precompiled binaries from the [Release page](https://github.com/tslj1024/Wormhole-Slim-UDP/releases)

### Source Code Compilation

Ensure Go 1.22+ is installed:

```bash
git clone https://github.com/tslj1024/Wormhole-Slim-UDP.git
cd Wormhole-Slim-UDP/server
go build -o server server.go
cd ../client
go build -o client client.go
```

## Quick Start

### Server (Public Server)

First, find a location on your public server and copy the `server` program and the `config` folder from the same folder as it, as shown below:

```
wormhole
├── server
└── config/
    └── app.yml
```

Modify the configuration file, `app.yml`:

```yaml
server:
  port: 8081 # This port is used to listen for client connections
  # This list is used to define the mapping of each public network server port to the intranet service
  clients:
    - clientId: CLIENTID1 # Differentiate each client. Copy first
      port: 8083  		  # The port through which users access intranet services
      tHost: 127.0.0.1 	  # IP address of the intranet service
      tPort: 80			  # Intranet service port
```

Run：

```bash
./server
```

### Client (Intranet Machine)

First, find a location on the intranet server and copy the `client` program and the `config` folder from the same folder as it, as shown below:

```
wormhole
├── client
└── config/
    └── app.yml
```

Modify the configuration file, `app.yml`:

```yaml
client:
  host: localhost     # Public server address
  port: 8081		  # Public network server port, the port used by the public network server to listen for client connections
  clientId: CLIENTID1 # Copy the CLIENTID copied from the server here
```

Run：

```bash
./client
```

## Security Recommendations

1. Do not disclose any of your Client IDs.
2. Use end-to-end encryption. This software does not provide encryption, and encryption is not required even with end-to-end encryption.
3. Expose only necessary services.

## FAQ

**Q: What should I do if the connection fails?**

A: Check the following:

- [ ] Check your firewall, especially for public servers.
- [ ] Ensure the client and server are configured correctly, especially the Client ID.
- [ ] Ensure intranet services are available.
- [ ] Ensure the client can access intranet services.
- [ ] Ensure the client can connect to the server. If it cannot connect, the client will output an error message.
- [ ] Or you can build services in the `test` folder on both sides to test connectivity

## Contribution Guidelines

Issues and pull requests are welcome. Before submitting code, please ensure the following:

1. Pass basic tests.
2. Follow existing coding styles.
3. Update relevant documentation.

## License

This project is licensed under the GPL-3.0 open source license. See the [LICENSE](LICENSE) file for details.

