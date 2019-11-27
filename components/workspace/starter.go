package workspace

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/config"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/starter"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tagger"
)

func Starter() starter.Operator {
	return &workspaceStarter{}
}

var l logger.Operator
var _ starter.Operator = &workspaceStarter{}

type workspaceStarter struct {
	interfaceKey joiner.InterfaceKey
}

func (ws *workspaceStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ws *workspaceStarter) Init(_ *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	ws.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil, nil
}

func (ws *workspaceStarter) Setup() error {
	return nil
}

func (ws *workspaceStarter) Run(joinerOp joiner.Operator) error {
	dataOp, ok := joinerOp.Interface(data.InterfaceKey).(data.Operator)
	if !ok {
		return errors.Errorf("no data.Operator with key %s", data.InterfaceKey)
	}

	taggerOp, ok := joinerOp.Interface(tagger.InterfaceKey).(tagger.Operator)
	if !ok {
		return errors.Errorf("no tagger.Operator with key %s", tagger.InterfaceKey)
	}

	wsOp, _, err := NewWorkspace(dataOp, taggerOp)
	if err != nil {
		return errors.Wrap(err, "can't init workspace.Operator")
	}

	err = joinerOp.Join(wsOp, ws.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *ws as workspace.Operator with key '%s'", ws.interfaceKey)
	}

	return nil
}
