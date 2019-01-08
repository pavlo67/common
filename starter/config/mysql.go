package config

type Table struct {
	Key  string
	Name string
}

type MySQLField struct {
	Name    string
	Type    string
	Null    bool
	Default string
	Extra   string
}

type MySQLIndex struct {
	Name      string   `json:"name"`
	Type      string   `json:"type,omitempty"`
	Fields    []string `json:"fields"`
	IndexType string   `json:"-"`
}

type MySQLTableComponent struct {
	Fields  [][]string   `json:"fields"`
	Indexes []MySQLIndex `json:"indexes"`
}
