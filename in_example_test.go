package use

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExampleIn(t *testing.T) {

	type Dest struct {
		F1 string
		F2 *int
		F3 *int
		F4 bool
	}

	// When we want to use In method, we define the tags on source struct.
	// This makes sense when we can have multiple different input sources
	// for one destination struct
	type Src struct {
		F1    *string `usein:""`
		F2inp int     `usein:"F2,nooverwrite"`
		F3inp int     `usein:"F3,nooverwrite"`
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

	setFields, err := In(&dest, &src)
	require.NoError(t, err)

	expectedDest := Dest{
		F1: "new f1",  // new value used from source
		F2: asRef(42), // will not be overwritten
		F3: asRef(44), // will be overwritten because dest original is nil
		F4: true,      // will not report error, even missing in source (omitmissing)
	}
	require.Equal(t, expectedDest, dest)

	expectedSetFields := []string{"F1", "F3"}
	sort.Strings(setFields)
	require.Equal(t, expectedSetFields, setFields)
}
