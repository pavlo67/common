package config

type SQLField struct {
	Name    string
	Type    string
	Null    bool
	Default string
	Extra   string
}

type SQLIndex struct {
	Name   string   `json:"name"`
	Type   string   `json:"type,omitempty"`
	Fields []string `json:"fields"`
}

type SQLTable struct {
	Name      string     `json:"name"`
	FieldsArr [][]string `json:"fields"`
	Fields    []SQLField `json:"-"`
	Indexes   []SQLIndex `json:"indexes"`
}
