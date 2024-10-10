package use

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFromSameName(t *testing.T) {
	type T struct {
		F1a   string         `usefrom:""`    // same name
		F1b   float64        `usefrom:"F1b"` // explicit (not needed, but testing it)
		F2a   string         `usefrom:",nooverwrite"`
		F2b   float64        `usefrom:"F2b,nooverwrite"`
		F3    *int           `usefrom:""`
		F3a   *int           `usefrom:""` // for sake of testing nil in output
		F3b   *int           `usefrom:""` // for sake of testing nil in output and input
		F4    int            `usefrom:""`
		F4nil int            `usefrom:""`
		F5    []int          `usefrom:""`
		F5a   []int          `usefrom:""` // for sake of testing nil input
		F7    map[string]int `usefrom:""`
		F7a   map[string]int `usefrom:""` // for sake of testing nil input
		Other int
	}

	obj := T{
		F1a:   "old f1a",                      // will be overwritten
		F1b:   99,                             // will be overwritten
		F2a:   "old f2a",                      // will not be overwritten
		F2b:   199,                            // will not be overwritten
		F3:    asRef(1234),                    // will be overwritten
		F3a:   nil,                            // will be overwritten
		F3b:   nil,                            //
		F4:    5678,                           // will be overwritten
		F4nil: 2121,                           // will not be overwritten
		F5:    []int{1, 2, 3},                 // will be overwritten
		F5a:   []int{1, 2, 3},                 // will not be overwritten
		F7:    map[string]int{"a": 1, "b": 2}, // will be overwritten
		F7a:   map[string]int{"a": 1, "b": 2}, // will not be overwritten
		Other: 123,
	}

	type TInput struct {
		F1a          string
		F1b          float64
		F2a          string
		F2b          float64
		F3           *int
		F3a          *int
		F3b          *int
		F4           *int
		F4nil        *int
		F5           []int
		F5a          []int
		F7           map[string]int
		F7a          map[string]int
		NotImportant string
	}

	objInput := TInput{
		F1a:          "some f1",
		F1b:          44,
		F2a:          "some f2",
		F2b:          55,
		F3:           asRef(422),
		F3a:          asRef(432),
		F3b:          nil,
		F4:           asRef(12345),
		F4nil:        nil, // explicit for readability
		F5:           []int{4, 5, 6},
		F5a:          nil, // explicit for readability
		F7:           map[string]int{"c": 3, "d": 4},
		F7a:          nil, // explicit for readability
		NotImportant: "not important",
	}
	objExpected := T{
		F1a:   "some f1", // used from input
		F1b:   44,        // used from input
		F2a:   "old f2a", // nooverwrite
		F2b:   199,       // nooverwrite
		F3:    asRef(422),
		F3a:   asRef(432),
		F3b:   nil, // explicit for readability
		F4:    12345,
		F4nil: 2121, // not used from input (nil)
		F5:    []int{4, 5, 6},
		F5a:   []int{1, 2, 3}, // not used from input (nil)
		F7:    map[string]int{"c": 3, "d": 4},
		F7a:   map[string]int{"a": 1, "b": 2}, // not used from input (nil)
		Other: 123,                            // not used from input (missing tag)
	}

	expectedSetFields := []string{"F1a", "F1b", "F3", "F3a", "F4", "F5", "F7"}
	sort.Strings(expectedSetFields)

	setfields, err := From(&obj, &objInput)
	sort.Strings(setfields)

	require.NoError(t, err)
	require.Equal(t, objExpected, obj)
	require.Equal(t, expectedSetFields, setfields)
}

func TestNewFromSameNameSameType(t *testing.T) {
	type T struct {
		F1a string  `usefrom:""`    // same name
		F1b float64 `usefrom:"F1b"` // explicit (not needed, but testing it)
		F2a string  `usefrom:",nooverwrite"`
		F2b float64 `usefrom:"F2b,nooverwrite"`
		F3  int
	}

	obj := T{
		F1a: "old f1a", // will be overwritten
		F1b: 99,        // will be overwritten
		F2a: "old f2a", // will not be overwritten
		F2b: 199,       // will not be overwritten
		F3:  123,
	}
	objInput := T{
		F1a: "some f1",
		F1b: 44,
		F2a: "some f2",
		F2b: 55,
		F3:  42,
	}
	objExpected := T{
		F1a: "some f1", // used from input
		F1b: 44,        // used from input
		F2a: "old f2a", // nooverwrite
		F2b: 199,       // nooverwrite
		F3:  123,       // not used from input (missing tag)
	}

	_, err := From(&obj, &objInput)

	require.NoError(t, err)
	require.Equal(t, objExpected, obj)
}

func TestNewFromMappedName(t *testing.T) {
	type T struct {
		F1a string  `usefrom:"InpF1a"`
		F1b string  `usefrom:"InpF1a"` // testing double mapping
		F2a string  `usefrom:"InpF2a,nooverwrite"`
		F2b float64 `usefrom:"InpF2b,nooverwrite"`
		F3  int
	}

	obj := T{
		F1a: "old f1a", // will be overwritten
		F1b: "old f1b", // will be overwritten
		F2a: "old f2a", // will not be overwritten
		F2b: 199,       // will not be overwritten
		F3:  123,
	}

	type TInput struct {
		InpF1a       string
		InpF2a       string
		InpF2b       float64
		F3           int
		NotImportant string
	}

	objInput := TInput{
		InpF1a: "some f1",
		InpF2a: "some f2",
		InpF2b: 55,
		F3:     42,
	}
	objExpected := T{
		F1a: "some f1", // used from input
		F1b: "some f1", // used from input
		F2a: "old f2a", // nooverwrite
		F2b: 199,       // nooverwrite
		F3:  123,       // not used from input (missing tag)
	}

	_, err := From(&obj, &objInput)

	require.NoError(t, err)
	require.Equal(t, objExpected, obj)
}

func TestNewFromOmitmissing(t *testing.T) {
	type T struct {
		F1a  string  `usefrom:",omitmissing"`
		F1aa string  `usefrom:",omitmissing"`
		F1b  float64 `usefrom:"F1b,omitmissing"`
		F1bb float64 `usefrom:"F1bb,omitmissing"`
		F2a  string  `usefrom:"NewF2a,nooverwrite,omitmissing"`
		F2b  float64 `usefrom:"NewF2b,nooverwrite,omitmissing"`
		F3   int
	}

	obj := T{
		F1a:  "old f1a",  // will be overwritten
		F1aa: "old f1aa", // will not be overwritten (no field in input)
		F1b:  99,         // will be overwritten
		F1bb: 999,        // will not be overwritten (no field in input)
		F2a:  "old f2a",  // will not be overwritten
		F2b:  199,        // will not be overwritten
		F3:   123,
	}

	type TInput struct {
		F1a          string
		F1b          float64
		NewF2a       string
		F2b          float64
		F3           int
		NotImportant string
	}

	objInput := TInput{
		F1a:          "some f1",
		F1b:          44,
		NewF2a:       "some f2",
		F2b:          55,
		F3:           42,
		NotImportant: "not important",
	}

	objExpected := T{
		F1a:  "some f1", // used from input
		F1aa: "old f1aa",
		F1b:  44, // used from input
		F1bb: 999,
		F2a:  "old f2a", // nooverwrite
		F2b:  199,       // same value (missing in input)
		F3:   123,       // not used from input (missing tag)
	}

	_, err := From(&obj, &objInput)

	require.NoError(t, err)
	require.Equal(t, objExpected, obj)
}

func TestNewFromPointer(t *testing.T) {
	type T struct {
		F1 string `usefrom:""` // same name
		F2 string `usefrom:",nooverwrite"`
		F3 int
		F4 int
	}

	obj := T{
		F1: "old f1a", // will be overwritten
		F2: "old f2a", // will not be overwritten
		F3: 123,
		F4: 42,
	}

	type TInput struct {
		F1           *string
		F2           *string
		F3           *int
		F4           int
		NotImportant string
	}

	objInput := TInput{
		F1:           asRef("some f1"),
		F2:           asRef("some f2"),
		F3:           nil, // to be explicit
		F4:           42,
		NotImportant: "not important",
	}
	objExpected := T{
		F1: "some f1", // used from input
		F2: "old f2a", // nooverwrite
		F3: 123,       // input has nil value
		F4: 42,        // not used from input (missing tag)
	}

	_, err := From(&obj, &objInput)

	require.NoError(t, err)
	require.Equal(t, objExpected, obj)
}

func TestNewFromNestedStruct(t *testing.T) {
	type Nested struct {
		NestedF1 string `usefrom:""`
	}

	type T struct {
		F1  string  `usefrom:""`
		F2  *Nested `usefrom:""`
		F3  *Nested `usefrom:",nooverwrite"`
		F4  *Nested `usefrom:""`
		F5  *Nested `usefrom:""`
		F6  Nested  `usefrom:""` // out direct, input pointer
		F7  Nested  `usefrom:""` // out direct, input pointer (nil)
		F8  Nested  `usefrom:""`
		F9  Nested  `usefrom:""`
		F10 Nested  `usefrom:""`
	}

	obj := T{
		F1:  "old f1a", // will be overwritten
		F2:  &Nested{NestedF1: "old nested f1"},
		F3:  &Nested{NestedF1: "old nested f1"},
		F4:  nil,
		F5:  nil,
		F6:  Nested{NestedF1: "old nested f1"},
		F7:  Nested{NestedF1: "old nested f1"},
		F8:  Nested{},
		F9:  Nested{},
		F10: Nested{NestedF1: "old nested f1"},
	}

	type TInput struct {
		F1           *string
		F2           *Nested
		F3           *Nested
		F4           *Nested
		F5           *Nested
		F6           *Nested
		F7           *Nested
		F8           *Nested
		F9           *Nested
		F10          Nested
		NotImportant string
	}

	objInput := TInput{
		F1:           asRef("new f1"),
		F2:           &Nested{NestedF1: "new nested f1"},
		F3:           &Nested{NestedF1: "new nested f1"},
		F4:           &Nested{NestedF1: "new nested f1"},
		F5:           nil,
		F6:           &Nested{NestedF1: "new nested f1"},
		F7:           nil,
		F8:           &Nested{NestedF1: "new nested f1"},
		F9:           nil,
		F10:          Nested{NestedF1: "new nested f1"},
		NotImportant: "not important",
	}
	objExpected := T{
		F1:  "new f1",
		F2:  &Nested{NestedF1: "new nested f1"},
		F3:  &Nested{NestedF1: "old nested f1"},
		F4:  &Nested{NestedF1: "new nested f1"},
		F5:  nil,
		F6:  Nested{NestedF1: "new nested f1"},
		F7:  Nested{NestedF1: "old nested f1"},
		F8:  Nested{NestedF1: "new nested f1"},
		F9:  Nested{},
		F10: Nested{NestedF1: "new nested f1"},
	}

	_, err := From(&obj, &objInput)

	require.NoError(t, err)
	require.Equal(t, objExpected, obj)
}

// ====

func TestNewFromUseFromErrors(t *testing.T) {
	type T struct {
		F1 string `usefrom:"NewF1"`
		F2 int    `usefrom:""`
		F6 string `usefrom:"NewF6,omitmissing"`
		F7 string `usefrom:",omitmissing,nooverwrite"`
	}

	t.Run("wrong type", func(t *testing.T) {
		type TInput struct {
			NewF1 string
			F2    bool // wrong type
			NewF6 string
			F7    string
		}

		obj := T{}
		objInput := TInput{
			NewF1: "some f1",
			F2:    true,
			NewF6: "some f6",
			F7:    "some f7",
		}

		_, err := From(&obj, &objInput)

		require.Error(t, err)
	})

	t.Run("missing field", func(t *testing.T) {
		type TInput struct {
			// NewF1 string // missing NewF1 and not omitmissing
			F2    int
			NewF6 string
			F7    string
		}

		obj := T{}
		objInput := TInput{
			F2:    42,
			NewF6: "some f6",
			F7:    "some f7",
		}

		_, err := From(&obj, &objInput)

		require.Error(t, err)
	})

	t.Run("input not a struct", func(t *testing.T) {
		obj := T{}
		objInput := 42

		_, err := From(&obj, &objInput)

		require.Error(t, err)
	})
}

func TestNewFromUseFromErrorsNotSettable(t *testing.T) {
	type T struct {
		// nolint here because I need to have unexported field (and try to use/set it)
		f1 string `usefrom:"NewF1"` //nolint:unused
		F2 int    `usefrom:""`
		F6 string `usefrom:"NewF6,omitmissing"`
		F7 string `usefrom:",omitmissing,nooverwrite"`
	}

	t.Run("not settable", func(t *testing.T) {
		type TInput struct {
			NewF1 string
			F2    int
			NewF6 string
			F7    string
		}

		obj := T{}
		objInput := TInput{
			NewF1: "some f1",
			F2:    42,
			NewF6: "some f6",
			F7:    "some f7",
		}

		_, err := From(&obj, &objInput)

		require.Error(t, err)
	})
}

// func TestFromUseErrorsNestedStruct(t *testing.T) {
// 	type Nested struct {
// 		F1 string `usefrom:""`
// 	}

// 	type T struct {
// 		F2 int     `usefrom:""`
// 		F6 *Nested `usefrom:""`
// 	}

// 	type TInput struct {
// 		F2 int
// 		F6 *Nested
// 	}

// 	obj := T{}
// 	objInput := TInput{
// 		F2: 42,
// 		F6: &Nested{
// 			F1: "some f1",
// 		},
// 	}

// 	err := From(&obj, &objInput)
// 	fmt.Printf("errrrrrrr: %v\n", err)

// 	require.Error(t, err)
// 	t.Fail()
// }