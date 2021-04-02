package def

type Flusher interface {
	// Flush sends any buffered data to the client.
	Flush()
}

type Closer interface {
	Close()
}
