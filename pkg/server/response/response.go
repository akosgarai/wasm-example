package response

// Socket holds the response data of the socket messages
type Socket struct {
	Error interface{}
	Data  map[string]string
}

// NewSocket returns a new Socket instance
func NewSocket() *Socket {
	return &Socket{
		Data: make(map[string]string),
	}
}
