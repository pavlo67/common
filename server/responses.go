package server

type BinaryResponse struct {
	Status   int
	Data     []byte
	MIMEType string
	FileName string
}

type DataResponse struct {
	Status int         `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}
