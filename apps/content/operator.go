package content

import "github.com/pavlo67/constructor/basis"

type Brief struct {
	ID      basis.ID   `bson:"_id,omitempty"     json:"id,omitempty"`
	Title   string     `bson:"title"             json:"title"`
	Summary string     `bson:"summary,omitempty" json:"summary,omitempty"`
	Info    basis.Info `bson:"info,omitempty"    json:"info,omitempty"`
}
