package util

// This file manages the mapping between the programming data types
// and the way those types are displayed in the user interface.
//

import (
	"fmt"
	"reflect"
	"strings"
)

// W2uiHTMLSelect is a struct that covers the way w2ui sends back the
// selection from a dropdown list.
type W2uiHTMLSelect struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// W2uiHTMLIdTextSelect is a struct that covers the way w2ui sends back the
// selection from a dropdown list.
type W2uiHTMLIdTextSelect struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// Str2Int64Map is a generic type for mapping strings and int64s
type Str2Int64Map map[string]int64

// assignmap defines the known type conversions. The mapper function can
// be called with reflect values for the two variables to map and the
// migration will be performed. Many of the conversions are between
// a list of strings and an int64.  For these conversions you can
// supply the Str2Int64Map and use the generic MigrateStrToInt64
// or MigrateInt64ToStr.
var assignmap = []struct {
	a      string                                           // source value
	b      string                                           // destination value
	mapper func(a, b *reflect.Value, m *Str2Int64Map) error // mapping function
	valmap *Str2Int64Map                                    // string to int64 map
}{
	{a: "int", b: "JSONbool", mapper: Int2Bool},
	{a: "JSONbool", b: "int", mapper: Bool2Int},
	{a: "int64", b: "JSONbool", mapper: Int642Bool},
	{a: "JSONbool", b: "int64", mapper: Bool2Int64},
	{a: "XJSONYesNo", b: "int", mapper: MigrateStrToInt64, valmap: &YesNoMap},
	{a: "int", b: "XJSONYesNo", mapper: MigrateInt64ToString, valmap: &YesNoMap},
}

var xjson = string("XJSON")

// XJSONprocess attempts to map a to b. If no converter can befound
// a message will be printed, then it will panic!
func XJSONprocess(a, b *reflect.Value) error {
	at := (*a).Type().String()
	bt := (*b).Type().String()
	// fmt.Printf("XJSONprocess: map from %s to %s\n", at, bt)
	for i := 0; i < len(assignmap); i++ {
		if strings.Contains(at, assignmap[i].a) && strings.Contains(bt, assignmap[i].b) {
			assignmap[i].mapper(a, b, assignmap[i].valmap)
			return nil
		}
	}
	return fmt.Errorf("XJSONmap - no conversion between: %s and %s", at, bt)
}

// MigrateStrToInt64 generic map of string to int64
func MigrateStrToInt64(a, b *reflect.Value, m *Str2Int64Map) error {
	si := (*a).Interface()
	s := fmt.Sprintf("%v", si)
	id, ok := (*m)[s]
	if !ok {
		id = int64(0)
	}
	(*b).Set(reflect.ValueOf(id))
	return nil
}

// MigrateInt64ToString generic mapping from int64 to enumerated strings
func MigrateInt64ToString(a, b *reflect.Value, m *Str2Int64Map) error {

	// fmt.Printf("Convert %d to BUD\n", (*a).Interface().(int64))

	s, err := (*m).ReverseMap((*a).Interface().(int64))
	if err != nil {
		return err
	}

	(*b).Set(reflect.ValueOf(s).Convert((*b).Type()))
	return nil
}

// XJSONYesNo is a UI converter: back-end int, UI: string
type XJSONYesNo string

// YesNoMap is the mapping for no = 0, 1 = yes
var YesNoMap = Str2Int64Map{
	"no":  int64(0),
	"yes": int64(1),
}

// Int2Bool copies an int into a bool value as follows
// if the int is 0, the bool value is false
// for any other value of the int the bool is true
// a must point to an int
// b must point to a bool
func Int2Bool(a, b *reflect.Value, m *Str2Int64Map) error {
	(*b).Set(reflect.ValueOf(0 != (*a).Interface().(int)))
	return nil
}

// Bool2Int is the exact inverse of Int2Bool
// a must point to a bool
// b must point to an int
func Bool2Int(a, b *reflect.Value, m *Str2Int64Map) error {
	i := 0
	if (*a).Interface().(bool) {
		i = 1
	}
	(*b).Set(reflect.ValueOf(i))
	return nil
}

// Int642Bool copies an int into a bool value as follows
// if the int is 0, the bool value is false
// for any other value of the int the bool is true
// a must point to an int
// b must point to a bool
func Int642Bool(a, b *reflect.Value, m *Str2Int64Map) error {
	(*b).Set(reflect.ValueOf(0 != (*a).Interface().(int64)))
	return nil
}

// Bool2Int64 is the exact inverse of Int642Bool
// a must point to a bool
// b must point to an int
func Bool2Int64(a, b *reflect.Value, m *Str2Int64Map) error {
	i := int64(0)
	if (*a).Interface().(bool) {
		i = int64(1)
	}
	(*b).Set(reflect.ValueOf(i))
	return nil
}

// RtActiveFLAGS is the mapping of FLAGS for active/inactive indication
var RtActiveFLAGS = Str2Int64Map{
	"Yes": int64(1),
	"No":  int64(0),
}
