package transport_http

import (
	"strconv"
	"strings"
	"time"

	"github.com/pavlo67/workshop/components/runner"

	"github.com/pkg/errors"

	"github.com/pavlo67/workshop/common/identity"
	"github.com/pavlo67/workshop/common/server/server_http"

	"github.com/pavlo67/workshop/components/packs"
	"github.com/pavlo67/workshop/components/transport"
	"github.com/pavlo67/workshop/components/transportrouter"
)

var _ transport.Operator = &transportHTTP{}

type transportHTTP struct {
	packsOp       packs.Operator
	runnerFactory runner.Factory
	routerOp      transportrouter.Operator

	domain identity.Domain
	path   string
	id     uint64

	routes transportrouter.Routes

	maxResponseDuration time.Duration

	//handlers map[[2]identity.Key]packs.Handler
	//mutex    *sync.RWMutex
}

// TODO!!! customize it
const maxResponseDuration = time.Second * 30

const onNew = "on sender_http.New(): "

func New(packsOp packs.Operator, runnerFactory runner.Factory, routerOp transportrouter.Operator, domain identity.Domain) (transport.Operator, *server_http.Endpoint, error) {
	if packsOp == nil {
		return nil, nil, errors.New(onNew + "no packs.Actor")
	}

	if runnerFactory == nil {
		return nil, nil, errors.New(onNew + "no runner.Factory")
	}

	if routerOp == nil {
		return nil, nil, errors.New(onNew + "no router.Actor")
	}

	if strings.TrimSpace(string(domain)) == "" {
		return nil, nil, errors.New("domain is empty")
	}

	routes, err := routerOp.Routes()
	if err != nil {
		// TODO: get routes later

		return nil, nil, errors.Wrap(err, onNew+"can't get routes")
	}

	//handlers := map[[2]identity.Key]packs.Handler{}

	transpOp := transportHTTP{
		packsOp:       packsOp,
		runnerFactory: runnerFactory,
		routerOp:      routerOp,

		routes: routes,
		domain: domain,
		path:   strconv.FormatInt(time.Now().UnixNano(), 10),

		maxResponseDuration: maxResponseDuration,
		//handlers: handlers,
		//mutex:    &sync.RWMutex{},
	}

	return &transpOp, transpOp.receiveEndpoint(), nil
}

// TODO!!! be careful because "pack.History = ..." isn't a thread safe action
