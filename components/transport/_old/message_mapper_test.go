package connector

import (
	"testing"

	"github.com/pavlo67/partes/crud"
)

func TestMessageMapper(t *testing.T) {
	crud.MapperTest(t, []crud.Mapper{&MessageMapper{}})
}
