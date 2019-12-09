package transport

import "time"

type DataType string

const DataItemsDataType DataType = "data_items"

type Packet struct {
	SourceURL string
	CreatedAt time.Time

	Type  DataType
	Data  interface{}
	MaxID uint64
}
