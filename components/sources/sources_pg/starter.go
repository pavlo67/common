package sources_pg

// import (
// 	"github.com/pkg/errors"

// 	"github.com/pavlo67/workshop/common"
// 	"github.com/pavlo67/workshop/common/config"
// 	"github.com/pavlo67/workshop/common/joiner"
// 	"github.com/pavlo67/workshop/common/logger"
// 	"github.com/pavlo67/workshop/common/starter"

// 	"github.com/pavlo67/workshop/components/packs"
// )

// func Starter() starter.Actor {
// 	return &packsPgStarter{}
// }

// var l logger.Actor
// var _ starter.Actor = &packsPgStarter{}

// type packsPgStarter struct {
// 	config       config.Access
// 	table        string
// 	interfaceKey joiner.InterfaceKey
// }

// func (ps *packsPgStarter) Name() string {
// 	return logger.GetCallInfo().PackageName
// }

// func (ps *packsPgStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Actor, options common.Map) ([]common.Map, error) {
// 	l = lCommon

// 	var cfgPg config.Access
// 	err := cfg.Value("pg", &cfgPg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ps.config = cfgPg
// 	ps.table, _ = options.String("table")
// 	ps.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(packs.InterfaceKey)))

// 	return nil, nil
// }

// func (ps *packsPgStarter) Setup() error {
// 	return nil
// }

// func (ps *packsPgStarter) Run(joinerOp joiner.Actor) error {
// 	packsOp, _, err := New(ps.config, ps.table, ps.interfaceKey)
// 	if err != nil {
// 		return errors.Wrap(err, "can't init packs.Actor")
// 	}

// 	err = joinerOp.Join(packsOp, ps.interfaceKey)
// 	if err != nil {
// 		return errors.Wrapf(err, "can't join &packsPg as packs.Actor with key '%s'", ps.interfaceKey)
// 	}

// 	return nil
// }
