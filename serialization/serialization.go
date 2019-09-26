package serialization

// Extension - Serialization
type Serialization interface {
	// Get content type unique id, recommended that custom implementations use values greater than 20.
	GetContentTypeId() int
	GetContentType() string
	Unmarshal(v interface{}) ([]byte, error)
	Marshal(data []byte, v interface{}) error
}
