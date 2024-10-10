package use

import "fmt"

// In copies values from src to dest. It uses tags on source struct to define the destination fields.
//
// Example:
//
//	type Dest struct {
//		F1 string
//		F2 *int
//		F3 *int
//		F4 bool
//	}
//
//	// When we want to use In method, we define the tags on source struct.
//	// This makes sense when we can have multiple different input sources
//	// for one destination struct
//	type Src struct {
//		F1    *string `usein:""`
//		F2inp int     `usein:"F2,nooverwrite"`
//		F3inp int     `usein:"F3,nooverwrite"`
//	}
//
//	dest := Dest{
//		F1: "original f1", // will be overwritten (there is no nooverwrite tag)
//		F2: asRef(42),     // will not be overwritten
//		F3: nil,           // will be overwritten because dest original is nil
//		F4: true,          // will not report error, even missing in source (omitmissing)
//	}
//
//	src := Src{
//		F1:    asRef("new f1"),
//		F2inp: 43,
//		F3inp: 44,
//	}
//
//	setFields, err := In(&dest, &src)
//
//	// dest is now:
//	//	F1: "new f1",  // new value used from source
//	//	F2: asRef(42), // will not be overwritten
//	//	F3: asRef(44), // will be overwritten because dest original is nil
//	//	F4: true,      // will not report error, even missing in source (omitmissing)
func In(dest, src any) (setFields []string, err error) {
	return in(dest, src, "")
}

// shadowed for not clogging API with optional argument
func in(dest, src any, parentFieldName string) (setFields []string, err error) {
	destObj, err := newObj(dest)
	if err != nil {
		return nil, fmt.Errorf("invalid value of destination object: %w (on path: %q)", err, parentFieldName)
	}

	srcObj, err := newObj(src)
	if err != nil {
		return nil, fmt.Errorf("invalid value of source object: %w (on path: %q)", err, parentFieldName)
	}

	for srcFieldName, tf := range srcObj.inTaggedFields {
		destFieldName := tf.tag.fieldName
		destT, ok := destObj.fieldType(destFieldName)
		if !ok {
			if tf.tag.omitMissing {
				continue
			}
			return nil, fmt.Errorf("destination field '%s%s' does not exist", parentFieldName, destFieldName)
		}

		isRecursive := containsStructOrPtrToStruct(destT)
		if !isRecursive {
			srcVal, _ := srcObj.field(srcFieldName)
			if isNil(srcVal) {
				continue
			}

			wasSet, err := destObj.setField(destFieldName, srcVal, tf.tag)
			if err != nil {
				return nil, fmt.Errorf("failed to set field %s%q: %w", parentFieldName, destFieldName, err)
			}
			if wasSet {
				setFields = append(setFields, addToFields(parentFieldName, destFieldName))
			}
			continue
		}

		// we have sub structs
		newsrc, exists, isnil, _ := srcObj.fieldRefAny(srcFieldName)
		if !exists {
			continue
		}
		if isnil {
			continue
		}

		newdest, _, isnil, _ := destObj.fieldRefAny(tf.tag.fieldName)
		if !isnil && tf.tag.noOverwrite {
			continue
		}

		// if isnil, we need to create reference to it and save it to the obj
		if isnil {
			err := destObj.createEmpty(destFieldName)
			if err != nil {
				return nil, fmt.Errorf("creating empty value: %w, (path: %q)", err, parentFieldName)
			}
			newdest, _, _, _ = destObj.fieldRefAny(destFieldName)
		}

		subParentFieldName := destFieldName
		if parentFieldName != "" {
			subParentFieldName = parentFieldName + "." + destFieldName
		}
		newSetFields, err := in(newdest, newsrc, subParentFieldName)
		if err != nil {
			return nil, err
		}

		setFields = append(setFields, newSetFields...)
	}

	return setFields, nil
}
