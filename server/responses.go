package server

type BinaryResponse struct {
	Status   int
	Data     []byte
	MIMEType string
	FileName string
}
