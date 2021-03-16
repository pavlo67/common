package transformer_persons_operator_pack

import (
	"fmt"
	"testing"
	"time"

	"github.com/pavlo67/common/common/db/db_sqlite"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/persons"
	"github.com/pavlo67/common/common/persons/persons_sqlite"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data_exchange/components/structures"
	"github.com/pavlo67/data_exchange/components/transformer"
	"github.com/pavlo67/data_exchange/components/transformer/transformer_test_scenarios"
)

func TestTransformPersonsOperatorPack(t *testing.T) {
	return

	_, cfgService, l := apps.PrepareTests(t, "../../../apps/_environments/", "test", "")

	components := []starter.Starter{
		{db_sqlite.Starter(), nil},
		{persons_sqlite.Starter(), nil},
		{Starter(), common.Map{"persons_cleaner_key": persons.InterfaceCleanerKey}},
	}

	label := "PERSONS_OPERATOR_PACK/TEST BUILD"
	joinerOp, err := starter.Run(components, cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	personsOp, _ := joinerOp.Interface(persons.InterfaceKey).(persons.Operator)
	require.NotNil(t, personsOp)

	transformOp, _ := joinerOp.Interface(InterfaceKey).(transformer.Operator)
	require.NotNil(t, transformOp)

	dataInitial := persons.Pack{
		PackDescription: structures.PackDescription{
			Fields: structures.Fields{},
			ItemDescription: structures.ItemDescription{
				Label: "title",
				// ErrorsMap: nil,
				// History:   nil,
				CreatedAt: time.Now(),
				// UpdatedAt: nil,
			},
		},
		Items: []persons.Item{
			{
				Identity: auth.Identity{
					Nickname: "wqerwqer",
					Roles:    nil,
				},
				ItemDescription: structures.ItemDescription{
					Info: common.Map{"xxx": "yyy", "zzz": 777},
					URN:  "urn1",
					// History:   nil,
					CreatedAt: time.Now(),
					// UpdatedAt: nil,
				},
			},
			{
				Identity: auth.Identity{
					Nickname: "wqerwqer2",
					Roles:    rbac.Roles{rbac.RoleUser},
				},
				ItemDescription: structures.ItemDescription{
					Info: common.Map{"xxx2": "yyy", "zzz2": 222},
					URN:  "urn2",
					// History:   nil,
					CreatedAt: time.Now(),
					// UpdatedAt: nil,
				},
			},
		},
	}

	dataInitial.Items[0].SetCreds(auth.Creds{auth.CredsEmail: "aaa@bbb.ccc"})
	dataInitial.Items[1].SetCreds(auth.Creds{auth.CredsEmail: "aaa2@bbb.ccc"})

	var params common.Map

	copyFinal, statFinal, dataFinal := transformer_test_scenarios.TestOperator(t, transformOp, params, dataInitial, true)

	//copyFinal, _ := transformOp.Copy(nil, params)
	t.Logf("COPY (INTERNAL) FINAL: %#v", copyFinal)

	//statFinal, _ := transformOp.Stat(nil, params)
	if statFinalStringer, ok := statFinal.(fmt.Stringer); ok {
		t.Logf("STAT (INTERNAL) FINAL: %s", statFinalStringer.String())
	} else {
		t.Logf("STAT (INTERNAL) FINAL: %#v", statFinal)
	}

	//dataFinal, _ := transformOp.Out(nil, params)
	t.Logf("DATA (OUT) FINAL: %#v", dataFinal)

}
