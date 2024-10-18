# GoProxli - A simple Go Proxy Server

A simple HTTP proxy server written in Go, designed to forward requests to a specified target server while injecting custom headers. This proxy can be used to simulate scenarios such as traffic redirection through API gateways or the addition of headers for testing purposes.

## Features

- **HTTP Proxy**: Forward incoming requests to a specified target server.
- **Custom Header Injection**: Modify or add headers to outgoing requests.
- **Cross-Platform Support**: Build and run on Linux, macOS, and Windows.
- **Dynamic Target Configuration**: Specify the target server via command-line arguments.
- **Extensible**: Easily extendable to add more features such as request/response logging, caching, etc.

## Use Cases

- **Simulate API Gateways**: Use the proxy to simulate traffic routing through an API gateway, where headers (such as authentication tokens or client information) are injected into each request.
- **Testing and Debugging**: Proxy requests to an API and inject headers like `X-Forwarded-For` or `Authorization` to test how the backend handles different headers.
- **Traffic Monitoring**: Modify the proxy to add logging or metrics to monitor traffic passing through the system.

## Getting Started

### Prerequisites

- Go 1.20+ installed on your machine.

### Installation

1. Clone this repository:

```sh
git clone https://github.com/osuritz/go-proxli.git
cd go-proxli
```

2. Build and run during development
```sh
go run ./src -target <server>
```

3. Build the Go binary:
```sh
go build -o go-proxli ./src/
```

4. Run the proxy server:
```sh
./go-proxli -target http://destination-server.com -port 8080
```

Replace `http://destination-server.com` with the actual server URL you want to proxy requests to. The proxy will listen on `localhost:8080` by default (you can change the port via the `-port` flag).

### Usage
```sh
./go-proxli -target http://your-target-server.com -port 8080
```
This command will start the proxy server on port `8080`, forwarding requests to the specified target server (`http://your-target-server.com`).

#### Injecting Custom Headers
To simulate traffic routing through an API gateway, you can modify the `ProxyHandler` function in `proxy.go` to inject custom headers. Here's an example of how you can add a header like `X-Proxy-By` to every request:
```go
// Add a custom header before forwarding the request
req.Header.Add("X-Proxy-By", "GoProxli Server")
```

Simply edit the code in `src/proxy.go` to add more headers as needed.


### Command-Line Flags
* `-target`: (Required) The URL of the target server to forward requests to.
* `-port`: (Optional) The port to run the proxy server on. Default is `8080`.

### Example
```sh
./go-proxli -target http://example.com -port 9090
```

This will start the proxy server on port `9090` and forward all requests to `http://example.com`.


## Development

### Modifying the Proxy
You can extend or modify the proxy server by editing the `src/proxy.go` file. For example, you can:
* Add more custom headers.
* Implement request/response logging.
* Add support for TLS/HTTPS connections.

### Testing
You can test the proxy server locally using `curl` or any HTTP client:
```sh
curl -I http://localhost:8080/some-endpoint
```

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

## Contributing
Feel free to open issues or submit pull requests. Contributions are always welcome!