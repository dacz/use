package use

import "fmt"

// From copies values from src to dest. It uses tags on destination struct to define the source fields.
//
// Example:
//
//	type Dest struct {
//		F1 string `usefrom:""`
//		F2 *int   `usefrom:"F2inp,nooverwrite"`
//		F3 *int   `usefrom:"F3inp,nooverwrite"`
//		F4 bool   `usefrom:"F4inp,omitmissing"`
//	}
//
//	type Src struct {
//		F1    *string
//		F2inp int
//		F3inp int
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
//	setFields, err := From(&dest, &src)
//
//	// dest is now:
//	//	F1: "new f1",  // new value used from source
//	//	F2: asRef(42), // will not be overwritten
//	//	F3: asRef(44), // will be overwritten because dest original is nil
//	//	F4: true,      // will not report error, even missing in source (omitmissing)
func From(dest, src any) (setFields []string, err error) {
	return from(dest, src, "")
}

// shadowed for not clogging API with optional argument
func from(dest, src any, parentFieldName string) (setFields []string, err error) {
	destObj, err := newObj(dest)
	if err != nil {
		return nil, fmt.Errorf("invalid value of destination object: %w (on path: %q)", err, parentFieldName)
	}

	srcObj, err := newObj(src)
	if err != nil {
		return nil, fmt.Errorf("invalid value of source object: %w (on path: %q)", err, parentFieldName)
	}

	for destFieldName, tf := range destObj.fromTaggedFields {
		destT := tf.sf.Type
		isRecursive := containsStructOrPtrToStruct(destT)
		if !isRecursive {
			srcVal, ok := srcObj.field(tf.tag.fieldName)
			if !ok {
				if tf.tag.omitMissing {
					continue
				}
				return nil, fmt.Errorf("invalid value of source field '%s%s'", parentFieldName, tf.tag.fieldName)
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
		newsrc, exists, isnil, _ := srcObj.fieldRefAny(tf.tag.fieldName)
		if !exists {
			continue
		}
		if isnil {
			continue
		}

		newdest, _, isnil, _ := destObj.fieldRefAny(destFieldName)
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
		newSetFields, err := from(newdest, newsrc, subParentFieldName)
		if err != nil {
			return nil, err
		}

		setFields = append(setFields, newSetFields...)
	}

	return setFields, nil
}
