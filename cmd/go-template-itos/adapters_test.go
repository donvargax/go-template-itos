package main

import (
	"reflect"
	"testing"
)

// TestEveryPortIsBoundOnDiskAndInProcess holds both providers of adapters to
// every port (decision 14): each field of kind interface, func, pointer, map
// or slice is bound by diskAdapters and by inProcessAdapters, so a port
// bound in one driver and not the other fails here. A new field fails it
// until both providers bind it.
func TestEveryPortIsBoundOnDiskAndInProcess(t *testing.T) {
	m := &memory{}
	if err := m.Setup(); err != nil {
		t.Fatal(err)
	}
	for name, provided := range map[string]adapters{"diskAdapters": diskAdapters(), "inProcessAdapters": m.inProcessAdapters()} {
		v := reflect.ValueOf(provided)
		for i := range v.NumField() {
			field := v.Type().Field(i)
			switch field.Type.Kind() {
			case reflect.Interface, reflect.Func, reflect.Pointer, reflect.Map, reflect.Slice:
				if v.Field(i).IsNil() {
					t.Errorf("%s leaves adapters.%s unbound: bind its port there", name, field.Name)
				}
			}
		}
	}
}
