# Compression

The framework automatically compresses HTTP responses based on the client's `Accept-Encoding` header.

## Supported Algorithms

| Algorithm | Header Value |
|-----------|-------------|
| gzip | `gzip` |
| deflate | `deflate` |

## Usage

No configuration needed — compression is automatic. The client requests compression via the `Accept-Encoding` header:

```bash
curl -H "Accept-Encoding: gzip" http://127.0.0.1:8080/api/data
```

The response includes:

```
Content-Encoding: gzip
```

## How It Works

The `CompressStd` interceptor runs at `Order = 1500` in the response phase:

1. Checks if `Accept-Encoding` contains a supported algorithm
2. If the response is a `RetAdapter` (Stream, Html, etc.):
   - Creates a pipe: original reader → compressor → pipe writer
   - Streams the compressed output via a new `Stream` response
   - Preserves original status code and headers
3. If the response is serialized JSON (`def.Content`):
   - Compresses the bytes in-memory
   - Wraps in a `Stream` response

## Implementation

```go
// Register custom compression algorithms
compress.CompressRegister["br"] = &brotliCompressor{}
```

The `Compress` interface:

```go
type Compress interface {
    New(io.Writer) io.WriteCloser
    ContentEncoding() string
}
```
