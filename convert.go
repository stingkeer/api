package api

/**
convert func result to []byte
*/
type Convert interface {
	convert(interface{}) []byte
}
