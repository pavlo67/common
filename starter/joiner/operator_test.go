package joiner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterface(t *testing.T) {
	joiner := New()

	const keyA1 ComponentKey = "KeyA1"
	structA1 := &StructA{}

	const keyA2 ComponentKey = "KeyA2"
	structA2 := &StructA{}

	joiner.JoinComponent(structA1, keyA1)
	joiner.JoinComponent(structA2, keyA2)

	structA1Joined, ok := joiner.Interface(keyA1).(InterfaceA)
	require.True(t, ok)
	require.Equal(t, structA1, structA1Joined)
}

func TestComponentsAll(t *testing.T) {
	joiner := New()

	const textA1 = "StructA.Action()"
	const keyA1 ComponentKey = "KeyA1"
	structA1 := &StructA{text: textA1}
	structA3 := &StructA{text: textA1}

	const keyA2 ComponentKey = "KeyA2"
	structA2 := &StructA{}

	const keyB1 ComponentKey = "KeyB1"
	structB1 := &StructB{}

	joiner.JoinComponent(structA1, keyA1)
	joiner.JoinComponent(structA3, keyA1)
	joiner.JoinComponent(structB1, keyB1)
	joiner.JoinComponent(structA2, keyA2)

	components := joiner.ComponentsAll(keyA1)
	require.Equal(t, 2, len(components))

	for _, component := range components {
		require.Equal(t, keyA1, component.Key)

		interfaceA, ok := component.Interface.(InterfaceA)
		require.True(t, ok)
		require.NotNil(t, interfaceA)

		text := interfaceA.ActionA()
		require.Equal(t, textA1, text)
	}

	require.Equal(t, 1, structA1.NumActionA)
	require.Equal(t, 1, structA3.NumActionA)
	require.Equal(t, 0, structA2.NumActionA)
}

func TestComponentsAllWithSignature(t *testing.T) {
	joiner := New()

	const textA1 = "StructA.Action()"
	const keyA1 ComponentKey = "KeyA1"
	structA1 := &StructA{text: textA1}
	structA3 := &StructA{text: textA1}

	const keyA2 ComponentKey = "KeyA2"
	structA2 := &StructA{text: textA1}

	const keyB1 ComponentKey = "KeyB1"
	structB1 := &StructB{}

	joiner.JoinComponent(structA1, keyA1)
	joiner.JoinComponent(structA3, keyA1)
	joiner.JoinComponent(structB1, keyB1)
	joiner.JoinComponent(structA2, keyA2)

	components := joiner.ComponentsAllWithInterface((*InterfaceA)(nil))
	require.Equal(t, 3, len(components))

	for _, component := range components {
		interfaceA, ok := component.Worker.(InterfaceA)
		require.True(t, ok)
		require.NotNil(t, interfaceA)

		text := interfaceA.ActionA()
		require.Equal(t, textA1, text)
	}

	require.Equal(t, 1, structA1.NumActionA)
	require.Equal(t, 1, structA3.NumActionA)
	require.Equal(t, 1, structA2.NumActionA)
}

func TestCloseAll(t *testing.T) {
	joiner := New()

	const textA1 = "StructA.Action()"
	const keyA1 ComponentKey = "KeyA1"
	structA1 := &StructA{text: textA1}
	structA3 := &StructA{text: textA1}

	const keyA2 ComponentKey = "KeyA2"
	structA2 := &StructA{text: textA1}

	const keyB1 ComponentKey = "KeyB1"
	structB1 := &StructB{}

	joiner.JoinComponent(structA1, keyA1)
	joiner.JoinComponent(structA3, keyA1)
	joiner.JoinComponent(structB1, keyB1)
	joiner.JoinComponent(structA2, keyA2)

	joiner.CloseAll()

	require.Equal(t, 1, structA1.NumClose)
	require.Equal(t, 1, structA2.NumClose)
	require.Equal(t, 1, structA3.NumClose)
	require.Equal(t, 1, structB1.NumClose)
}

// InterfaceA (includes Closer) -----------------------------------------------------------------------------------------------------

type InterfaceA interface {
	ActionA() string
}

type StructA struct {
	NumActionA, NumClose int
	text                 string
}

var _ InterfaceA = &StructA{}
var _ Closer = &StructA{}

func (s *StructA) ActionA() string {
	s.NumActionA++
	fmt.Println("StructA.Action()")
	return s.text
}

func (s *StructA) Close() error {
	s.NumClose++
	fmt.Println("StructA.Close()")
	return nil
}

// InterfaceB (includes Closer) -----------------------------------------------------------------------------------------------------

type InterfaceB interface {
	ActionB() string
}

type StructB struct {
	NumActionB, NumClose int
}

var _ InterfaceB = &StructB{}
var _ Closer = &StructB{}

func (s *StructB) ActionB() string {
	s.NumActionB++
	fmt.Println("StructB.Action()")
	return "StructB.Action()"
}

func (s *StructB) Close() error {
	s.NumClose++
	fmt.Println("StructB.Close()")
	return nil
}
