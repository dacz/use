package use

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Somenested struct {
	Na int
}

type Someobj struct {
	B     string          `usefrom:"" usein:""`
	Br    *string         `usefrom:"" usein:""`
	Brnil *string         `usefrom:"" usein:""`
	C     map[string]int  `usefrom:"" usein:""`
	Cr    *map[string]int `usefrom:"" usein:""`
	Crnil map[string]int  `usefrom:"" usein:""`
	D     []int           `usefrom:"" usein:""`
	Dr    *[]int          `usefrom:"" usein:""`
	Drnil []int           `usefrom:"" usein:""`
	E     Somenested      `usefrom:"" usein:""`
	Er    *Somenested     `usefrom:"" usein:""`
	Ernil *Somenested     `usefrom:"" usein:""`
	F     []Somenested    `usefrom:"" usein:""`
	Fr    []*Somenested   `usefrom:"" usein:""`
	Frnil []Somenested    `usefrom:"" usein:""`
	Nt    int
}

type SomeAddresableValues struct {
	StringF        string
	StringRef      *string
	StringRefNil   *string
	MapF           map[string]int
	MapRef         *map[string]int
	MapRefNil      *map[string]int
	MapNil         map[string]int
	SliceF         []int
	SliceRef       *[]int
	SliceRefNil    *[]int
	SliceNil       []int
	StructF        Somenested
	StructRef      *Somenested
	StructRefNil   *Somenested
	SliceOfStructF []Somenested
	SliceOfPtrF    []*Somenested
}

var SomeAddresableObj = SomeAddresableValues{
	StringF:        "newstring",
	StringRef:      asRef("newstringref"),
	StringRefNil:   nil,
	MapF:           map[string]int{"newmap": 123321},
	MapRef:         asRef(map[string]int{"newmapref": 123321}),
	MapRefNil:      nil,
	MapNil:         nil,
	SliceF:         []int{123, 321, 1},
	SliceRef:       asRef([]int{123, 321, 1, 987}),
	SliceRefNil:    nil,
	SliceNil:       nil,
	StructF:        Somenested{123213},
	StructRef:      asRef(Somenested{123678}),
	StructRefNil:   nil,
	SliceOfStructF: []Somenested{{123}, {321}, {1}, {2}},
	SliceOfPtrF:    []*Somenested{{1234}, {4321}, {41}, {42}},
}

func newtestobj(t *testing.T) Someobj {
	t.Helper()
	return Someobj{
		B:     "two",
		Br:    asRef("twor"),
		Brnil: nil,
		C:     map[string]int{"three": 3},
		Cr:    asRef(map[string]int{"threer": 3}),
		Crnil: nil,
		D:     []int{4, 5},
		Dr:    asRef([]int{44, 55}),
		Drnil: nil,
		E:     Somenested{6},
		Er:    asRef(Somenested{66}),
		Ernil: nil,
		F:     []Somenested{{123}, {456}},
		Fr:    []*Somenested{{123}, {456}},
		Frnil: nil,
		Nt:    989,
	}
}

func TestNewObj(t *testing.T) {
	ov := newtestobj(t)
	o, err := newObj(&ov)
	require.NoError(t, err)
	require.True(t, o.isIndirect)
	require.Equal(t, o.derefType(), reflect.TypeOf(ov))
	require.Equal(t, o.derefValue().Interface(), reflect.ValueOf(ov).Interface())

	v, ok := o.field("Brnil")
	require.True(t, ok)
	require.True(t, isNil(v))

	// ---

	// no reference
	o, err = newObj(ov)
	require.Error(t, err)
	require.Nil(t, o)

	// nil reference
	var objnil *Someobj
	ornil, err := newObj(objnil)
	require.Error(t, err)
	require.Nil(t, ornil)

}
