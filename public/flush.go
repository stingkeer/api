package public

type Flusher interface {
	// Flush sends any buffered data to the client.
	Flush()
}
