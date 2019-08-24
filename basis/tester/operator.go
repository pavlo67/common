package tester

// operator ------------------------------------------------------------------------------

type CallKey string

type Call func(interface{}) (interface{}, error)

type Operator interface {
	Caller() map[CallKey]Call
	Scenarios() []Scenario
}
