# package use

[![Go Reference](https://pkg.go.dev/badge/github.com/dacz/use.svg)](https://pkg.go.dev/github.com/dacz/use) ![Tests](https://github.com/dacz/use/actions/workflows/ci.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/dacz/use)](https://goreportcard.com/report/github.com/dacz/use)

## Install

`go get -u github.com/dacz/use`

## Description

Package `use` provides a function to copy values from one struct to another.
Both types (destination and source) has to be structs and copied values
has to have the same types (but might mi mix of direct values and pointers)
usual example is to have some object and you get input with optional values
and you want to update the object with these values. Not all fields must be used
and some fields may use different names.

Example usage: You have an object (eg. got from DB) and user submitted update input comes from API (with optional values). Values that are not nil should be copied to the object (and know which ones were copied).

The definition is done with struct tags.

There are generally two ways to use this package:
  - use.From to copy values from source to destination,
    where the copy rules are defined in the destination struct tags
  - use.In to copy values from source to destination,
    where the copy rules are defined in the source struct tags

## Tags

`usefrom` defined on destination struct

    F1 <type> `usefrom:""` // will use the same field name from source struct
    // and it is an error if the `F1` field is missing in source struct.
    F2 <type> `usefrom:"F2inp"` // will use the field `F2inp` from source struct
    F3 <type> `usefrom:",nooverwrite"` // will use the same field name from source struct
    // and if destination struct has non nil value, it will not be overwritten
    F4 <type> `usefrom:",omitmissing"` // will try use the same field name from source struct
    // but if the `F4` field is missing in source struct, it will NOT report a problem

`usein` defined on source struct

    F1 <type> `usein:""` // will update the same field name as in source struct (F1)
    // and it is an error if the `F1` field is missing in destination struct.
    ... and similar renaming, `nooverwrite` and `omitmissing` as in `usefrom`.

## Examples

See [From test example](from_example_test.go) or [In test example](in_example_test.go). There are more test files to consult for details, nested structs, handling nils, ...

For more examples of usage, see use_test.go.

## To do:
  - [ ] support for interfaces
  - [ ] source as map[string]any. Usually we have some kind of validator working directly on data from API/form and output of this validator is usually struct. So the map[string]any is not needed in most of the cases.
  - [ ] LOWPRIORITY type cache for faster processing
  - [ ] WONTFIX support for maps and slices with nested structs (and transforming them)
    (conversion of map to input struct should happen somewhere else)
