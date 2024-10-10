package use

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExampleFrom(t *testing.T) {

	// When we want to use From method, we define the tags on destination struct.
	type Dest struct {
		F1 string `usefrom:""`
		F2 *int   `usefrom:"F2inp,nooverwrite"`
		F3 *int   `usefrom:"F3inp,nooverwrite"`
		F4 bool   `usefrom:"F4inp,omitmissing"`
	}

	type Src struct {
		F1    *string
		F2inp int
		F3inp int
	}

	dest := Dest{
		F1: "original f1", // will be overwritten (there is no nooverwrite tag)
		F2: asRef(42),     // will not be overwritten
		F3: nil,           // will be overwritten because dest original is nil
		F4: true,          // will not report error, even missing in source (omitmissing)
	}

	src := Src{
		F1:    asRef("new f1"),
		F2inp: 43,
		F3inp: 44,
	}

	setFields, err := From(&dest, &src)
	require.NoError(t, err)

	expectedDest := Dest{
		F1: "new f1",  // new value used from source
		F2: asRef(42), // will not be overwritten
		F3: asRef(44), // will be overwritten because dest original is nil
		F4: true,      // will not report error, even missing in source (omitmissing)
	}
	require.Equal(t, expectedDest, dest)

	expectedSetFields := []string{"F1", "F3"}
	// slices.Sort(setFields) // I want to support go 1.19, where slices pkg is not available
	sort.Strings(setFields)
	sort.Strings(expectedSetFields)
	require.Equal(t, expectedSetFields, setFields)
}
