package protocol

type Request struct {
	Header Header
	Key    []byte
	Value  []byte
	TTL    uint32
}

type Response struct {
	Header Header
	Status Status
	Value  []byte
	Error  string
}
