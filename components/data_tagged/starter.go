package data_tagged

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
	dataKey      joiner.InterfaceKey
	interfaceKey joiner.InterfaceKey
	noTagger     bool
}

func (ws *workspaceStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ws *workspaceStarter) Init(_, _ *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	ws.dataKey = joiner.InterfaceKey(options.StringDefault("data_key", string(data.InterfaceKey)))
	ws.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))
	ws.noTagger, _ = options.Bool("no_tagger")

	return nil, nil
}

func (ws *workspaceStarter) Setup() error {
	return nil
}

func (ws *workspaceStarter) Run(joinerOp joiner.Operator) error {
	dataOp, ok := joinerOp.Interface(ws.dataKey).(data.Operator)
	if !ok {
		return errors.Errorf("no data.Operator with key %s", ws.dataKey)
	}

	var taggerOp tagger.Operator

	if !ws.noTagger {
		taggerOp, ok = joinerOp.Interface(tagger.InterfaceKey).(tagger.Operator)
		if !ok {
			return errors.Errorf("no tagger.Operator with key %s", tagger.InterfaceKey)
		}
	}

	wsOp, _, err := New(dataOp, taggerOp)
	if err != nil {
		return errors.Wrap(err, "can't init storage.Operator")
	}

	err = joinerOp.Join(wsOp, ws.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *ws as storage.Operator with key '%s'", ws.interfaceKey)
	}

	return nil
}
