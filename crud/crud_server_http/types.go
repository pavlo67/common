package crud_serverhttp_jschmhr0

//import (
//	"github.com/pavlo67/punctum/basis"
//	"github.com/pavlo67/punctum/basis/crud"
//	"github.com/pavlo67/punctum/basis/selectors"
//)
//
//type AdminCount struct {
//	FieldKey string `json:"field_key"`
//	Title    string `json:"title"`
//	Count    uint   `json:"count"`
//}
//
//type AdminField struct {
//	crud.Field
//	Title string `json:"title"`
//	Check string `json:"check"`
//	Sort  string `json:"sort"`
//}
//
//type AdminQuery struct {
//	Title    string       `json:"title"`
//	TableKey string       `json:"table_key"`
//	Fields   []AdminField `json:"fields"`
//	crud.ReadOptions
//}
//
//type AdminDescription struct {
//	Title           string                `json:"title"`
//	PrimaryQueryKey string                `json:"primary_query_key"`
//	Queries         map[string]AdminQuery `json:"queries"`
//}
//
//type Operator interface {
//	// Describe ...
//	Describe() (*AdminDescription, error)
//
//	// Count
//	Count(identity auth.IDentity, selector selectors.Selector, options *crud.ReadOptions) ([]AdminCount, error)
//}
