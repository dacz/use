package use

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type obj struct {
	t                reflect.Type
	isIndirect       bool
	fromTaggedFields map[string]*taggedField // key is field name
	inTaggedFields   map[string]*taggedField // key is field name
	v                reflect.Value
}

func newObj(ov any) (*obj, error) {
	t := reflect.TypeOf(ov)
	v := reflect.ValueOf(ov)

	if indirect, derefed := derefType(t); !indirect || derefed.Kind() != reflect.Struct {
		return nil, errors.New("obj must be a reference to struct")
	}

	if v.IsNil() {
		return nil, errors.New("obj must be a reference to non nil object")
	}

	o := &obj{
		t:          t,
		isIndirect: t.Kind() == reflect.Ptr,
		v:          v,
	}

	derefedt := o.derefType()
	derefedv := o.derefValue()

	for i := 0; i < derefedt.NumField(); i++ {
		tagfrom := parseTag(derefedt.Field(i), fromTag)
		if tagfrom != nil {
			vf := derefedv.Field(i)
			if !vf.CanSet() {
				return nil, fmt.Errorf("field %q is not settable", derefedt.Field(i).Name)
			}

			if o.fromTaggedFields == nil {
				o.fromTaggedFields = make(map[string]*taggedField, 0)
			}
			o.fromTaggedFields[derefedt.Field(i).Name] = &taggedField{
				sf:         derefedt.Field(i),
				isIndirect: derefedt.Field(i).Type.Kind() == reflect.Ptr,
				tag:        tagfrom,
			}
		}

		tagin := parseTag(derefedt.Field(i), inTag)
		if tagin != nil {
			if o.inTaggedFields == nil {
				o.inTaggedFields = make(map[string]*taggedField, 0)
			}
			o.inTaggedFields[derefedt.Field(i).Name] = &taggedField{
				sf:         derefedt.Field(i),
				isIndirect: derefedt.Field(i).Type.Kind() == reflect.Ptr,
				tag:        tagin,
			}
		}
	}

	return o, nil
}

func (o *obj) tag(tagk tagKind, fname string) *tag {
	if tagk == fromTag {
		if tf, ok := o.fromTaggedFields[fname]; ok {
			return tf.tag
		}
		return nil
	}

	if tf, ok := o.inTaggedFields[fname]; ok {
		return tf.tag
	}
	return nil
}

func (o *obj) derefType() reflect.Type {
	if o.isIndirect {
		return o.t.Elem()
	}
	return o.t
}

func (o *obj) derefValue() reflect.Value {
	return reflect.Indirect(o.v)
}

func (o *obj) field(fname string) (reflect.Value, bool) {
	fv := o.derefValue().FieldByName(fname)
	if !fv.IsValid() {
		return reflect.Value{}, false
	}
	return fv, true
}

func (o *obj) fieldType(fname string) (reflect.Type, bool) {
	sf, ok := o.derefType().FieldByName(fname)
	if !ok {
		return nil, false
	}
	return sf.Type, true
}

func (o *obj) fieldRefAny(fname string) (v any, exists bool, isNilV bool, indirect bool) {
	var fv reflect.Value
	fv, exists = o.field(fname)
	if !exists {
		return
	}

	isNilV = isNil(fv)
	if isNilV {
		return
	}

	indirect, _ = derefValue(fv)
	if indirect {
		v = fv.Interface()
		return
	}
	v = fv.Addr().Interface()
	return
}

// setField does NOT set field if source is nil
func (o *obj) setField(fname string, v reflect.Value, tg *tag) (wasSet bool, err error) {
	if isNil(v) {
		return false, nil
	}

	fv, ok := o.field(fname)
	if !ok {
		if tg.omitMissing {
			return false, nil
		}
		return false, errors.New("dest field not found")
	}

	if !fv.CanSet() {
		return false, errors.New("dest field not settable")
	}

	if tg.noOverwrite && !isNil(fv) {
		return false, nil
	}

	if !typesMatch(fv.Type(), v.Type()) {
		return false, fmt.Errorf("types not assignable. dest %q, src %q", fv.Type().Kind(), v.Type().Kind())
	}

	if fv.Kind() == reflect.Ptr {
		if v.Kind() == reflect.Ptr {
			fv.Set(v)
			return true, nil
		}
		if !v.CanAddr() {
			return false, errors.New("cannot take address of source value")
		}
		fv.Set(v.Addr())
		return true, nil
	}

	if v.Kind() == reflect.Ptr {
		fv.Set(v.Elem())
		return true, nil
	}

	fv.Set(v)

	return true, nil
}

func (o *obj) createEmpty(fieldname string) error {
	ft, ok := o.derefType().FieldByName(fieldname)
	if !ok {
		return fmt.Errorf("%q is not defined on struct", fieldname)
	}

	_, derefedT := derefType(ft.Type)

	tv := reflect.New(derefedT)

	tg := o.tag(fromTag, fieldname)
	if tg == nil {
		return fmt.Errorf("tag for field %q does not exist", fieldname) // should not happen
	}
	_, err := o.setField(fieldname, tv, tg)
	if err != nil {
		return fmt.Errorf("error initializing new empty struct value: %w", err)
	}

	return nil
}

type taggedField struct {
	sf         reflect.StructField // has Name, Type, Tag
	isIndirect bool
	tag        *tag
}

type tagKind string

const (
	inTag   tagKind = "usein"
	fromTag tagKind = "usefrom"
)

type tag struct {
	fieldName   string
	noOverwrite bool
	omitMissing bool
}

func parseTag(structField reflect.StructField, tagName tagKind) *tag {
	v, ok := structField.Tag.Lookup(string(tagName))
	if !ok {
		return nil
	}

	tg := &tag{}

	vparts := strings.Split(v, ",")
	if len(vparts) != 0 {
		tg.fieldName = vparts[0]
	}

	for _, vpart := range vparts[1:] {
		if vpart == "nooverwrite" {
			tg.noOverwrite = true
			continue
		}
		if vpart == "omitmissing" {
			tg.omitMissing = true
			continue
		}
	}

	// no renaming, use field name
	if tg.fieldName == "" {
		tg.fieldName = structField.Name
	}

	return tg
}
