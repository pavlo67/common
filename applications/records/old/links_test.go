package old

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/constructor/components/auth"
	"github.com/pavlo67/constructor/confidenter"
	"github.com/pavlo67/constructor/confidenter/groups/groupsstub"
	"github.com/pavlo67/constructor/notebook/notes"
)

func TestFilterLinks(t *testing.T) {

	is1 := auth.ID("aaa/user/1")
	is2 := auth.ID("aaa/user/2")
	is3 := auth.ID("aaa/user/3")

	identity1 := is1.Identity()
	identity2 := is2.Identity()
	identity3 := is3.Identity()

	err := confidenter.SetSystemIdentity(map[string]string{"domain": "aaa"})
	require.NoError(t, err)

	groupIDs := map[auth.ID][]string{
		is1: {"1", "5", "6"},
		is2: {"2", "6"},
	}
	ctrlOp, _ := groupsstub.New(groupIDs, "")

	links1 := []notes.Item{
		{ID: "1", Type: "2", Name: "3", To: "4", RView: is1, ROwner: ""},
		{ID: "1", Type: "2", Name: "3", To: "4", RView: "", ROwner: is1},
	}

	links2 := []notes.Item{
		{ID: "1", Type: "2", Name: "3", To: "4", RView: is2, ROwner: ""},
		{ID: "1", Type: "2", Name: "3", To: "4", RView: "", ROwner: is2},
	}

	links3 := []notes.Item{
		{ID: "1", Type: "2", Name: "3", To: "4", RView: is3, ROwner: ""},
		{ID: "1", Type: "2", Name: "3", To: "4", RView: "", ROwner: is3},
	}

	links12 := append(links1, links2...)
	links13 := append(links1, links3...)
	links22 := append(links2, links2...)
	links23 := append(links2, links3...)
	links123 := append(links12, links3...)
	links122 := append(links12, links2...)
	links132 := append(links13, links2...)

	require.EqualValues(t, []notes.Item(nil), FilterLinks(&identity1, ctrlOp, links23))
	require.EqualValues(t, []notes.Item(nil), FilterLinks(&identity2, ctrlOp, links13))
	require.EqualValues(t, []notes.Item(nil), FilterLinks(&identity3, ctrlOp, links12))
	require.EqualValues(t, []notes.Item(nil), FilterLinks(&identity3, ctrlOp, links122))

	require.EqualValues(t, links1, FilterLinks(&identity1, ctrlOp, links1))
	require.EqualValues(t, links1, FilterLinks(&identity1, ctrlOp, links12))
	require.EqualValues(t, links1, FilterLinks(&identity1, ctrlOp, links123))
	require.EqualValues(t, links1, FilterLinks(&identity1, ctrlOp, links122))
	require.EqualValues(t, links1, FilterLinks(&identity1, ctrlOp, links132))

	require.EqualValues(t, links2, FilterLinks(&identity2, ctrlOp, links12))
	require.EqualValues(t, links2, FilterLinks(&identity2, ctrlOp, links23))
	require.EqualValues(t, links2, FilterLinks(&identity2, ctrlOp, links123))
	require.EqualValues(t, links2, FilterLinks(&identity2, ctrlOp, links132))

	require.EqualValues(t, links3, FilterLinks(&identity3, ctrlOp, links13))
	require.EqualValues(t, links3, FilterLinks(&identity3, ctrlOp, links23))
	require.EqualValues(t, links3, FilterLinks(&identity3, ctrlOp, links132))
	require.EqualValues(t, links3, FilterLinks(&identity3, ctrlOp, links123))

	require.EqualValues(t, links22, FilterLinks(&identity2, ctrlOp, links22))
	require.EqualValues(t, links22, FilterLinks(&identity2, ctrlOp, links122))
}
