# HTTP/3 (QUIC)

The framework supports HTTP/3 via the QUIC protocol using the `quic-go` library.

## Building with HTTP/3

```bash
go build -tags "http3"
```

This compiles `http3.go` instead of `http.go`, replacing the standard HTTP server with a dual HTTP/1.1 + HTTP/3 server.

## Starting the Server

In HTTP/3 mode, only `StartTLSService` is available — QUIC requires TLS:

```go
api.StartTLSService(
    api.WithListen(":443"),
    api.WithTLSFile("cert.pem", "key.pem"),
)
```

`StartService()` returns an error in HTTP/3 mode: `"http3 only support tls"`.

## How It Works

When built with the `http3` tag, `StartTLSService` starts **two servers concurrently**:

1. **TCP + TLS server** — standard HTTPS on the configured address
2. **QUIC server** — HTTP/3 on the same address using UDP

The TCP server adds `Alt-Svc` response headers to advertise HTTP/3 availability:

```
Alt-Svc: h3=":443"; ma=2592000
```

This tells browsers they can connect via HTTP/3 for subsequent requests.

## Error Handling

The framework waits for either server to fail:

```go
select {
case err := <-hErr:
    quicServer.Close()
    return err
case err := <-qErr:
    return err
}
```

If either server fails, the other is shut down.

## Client Support

Most modern browsers support HTTP/3:
- Chrome (since v87)
- Firefox (since v88)
- Safari (since v16)
- Edge (since v87)

`curl` supports HTTP/3 with the `--http3` flag:

```bash
curl --http3 https://localhost:443/api/data
```
