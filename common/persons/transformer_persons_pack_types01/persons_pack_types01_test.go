package transformer_persons_pack_types01

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data_exchange/components/structures"
	"github.com/pavlo67/data_exchange/components/transformer"
	"github.com/pavlo67/data_exchange/components/transformer/transformer_test_scenarios"
	"github.com/pavlo67/data_exchange/components/types/types01"
)

func TestTransformTableCSV(t *testing.T) {
	_, cfgService, l := apps.PrepareTests(t, "../../../apps/_environments/", "test", "")

	components := []starter.Starter{
		{Starter(), nil},
	}

	label := "PACK_PERSONS_TYPES01/TEST BUILD"
	joinerOp, err := starter.Run(components, cfgService, label, l)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	transformOp, _ := joinerOp.Interface(InterfaceKey).(transformer.Operator)
	require.NotNil(t, transformOp)

	dataInitial := structures.Pack{
		PackDescription: structures.PackDescription{
			Title:  "title",
			Fields: structures.Fields{},
			// ErrorsMap: nil,
			// History:   nil,
			CreatedAt: time.Now(),
			// UpdatedAt: nil,
		},
		Data: []types01.Person{
			{
				URN:      "urn1",
				Nickname: "wqerwqer",
				Roles:    nil,
				Creds:    common.Map{string(auth.CredsEmail): "aaa@bbb.ccc"},
				Data:     common.Map{"xxx": "yyy", "zzz": 777},
				// History:   nil,
				CreatedAt: time.Now(),
				// UpdatedAt: nil,
			},
			{
				URN:      "urn2",
				Nickname: "wqerwqer2",
				Roles:    rbac.Roles{rbac.RoleUser},
				Creds:    common.Map{string(auth.CredsEmail): "aaa2@bbb.ccc"},
				Data:     common.Map{"xxx2": "yyy", "zzz2": 222},
				// History:   nil,
				CreatedAt: time.Now(),
				// UpdatedAt: nil,
			},
		},
	}

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
