// Package `use` provides a function to copy values from one struct to another.
// Both types (destination and source) has to be structs and copied values
// has to have the same types (but might mi mix of direct values and pointers)
//
// usual example is to have some object and you get input with optional values
// and you want to update the object with these values. Not all fields must be used
// and some fields may use different names.
//
// The definition is done with struct tags.
//
// There are generally two ways to use this package:
//   - use.From to copy values from source to destination,
//     where the copy rules are defined in the destination struct tags
//   - use.In to copy values from source to destination,
//     where the copy rules are defined in the source struct tags
//
// For examples of usage, see use_test.go.
//
// # To do:
//
//   - support for interfaces
//   - source as map[string]any
//   - LOW type cache for faster processing
//   - WONTFIX support for maps and slices with nested structs (and transforming them)
//     (conversion of map to input struct should happen somewhere else)
package use
