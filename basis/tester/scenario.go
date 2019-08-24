package tester

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/basis/joiner"
)

// test scenario -------------------------------------------------------------------------

type InterfaceSelector struct {
	InterfaceKey   joiner.InterfaceKey
	PtrToInterface interface{}
}

type Play func(*testing.T, ...Operator)

type Scenario struct {
	InterfaceSelectors []InterfaceSelector
	Play
}

func (scenario Scenario) PlayWithJoiner(t *testing.T, joinerOp joiner.Operator) {
	var ops []Operator

	// TODO: play on carthesian product of interface variants if InterfaceKey is empty

	for i, interfaceSelector := range scenario.InterfaceSelectors {
		intrfc := joinerOp.Interface(interfaceSelector.InterfaceKey)
		require.Truef(
			t, reflect.TypeOf(intrfc).Implements(reflect.TypeOf(interfaceSelector.PtrToInterface).Elem()),
			"interface #%i is %T and doesn't impement %s interface", i, intrfc, reflect.TypeOf(interfaceSelector.PtrToInterface),
		)

		op, ok := intrfc.(Operator)
		require.Truef(t, ok, "interface #%i is %T and doesn't impement tester.Operator interface", i, intrfc)

		ops = append(ops, op)
	}

	scenario.Play(t, ops...)
}

func (scenario Scenario) Apply(t *testing.T, intfcs ...interface{}) {
	var ops []Operator
	for i, intrfc := range intfcs {
		op, ok := intrfc.(Operator)
		require.Truef(t, ok, "interface #%i is %T and doesn't impement tester.Operator interface", i, intrfc)

		ops = append(ops, op)
	}

	scenario.Play(t, ops...)
}
