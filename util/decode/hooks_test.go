package decode

import (
	"reflect"
	"testing"
)

func TestComposeDecodeHookFunc(t *testing.T) {
	f1 := func(
		f reflect.Kind,
		t reflect.Kind,
		data interface{}) (interface{}, error) {
		return data.(string) + "monkey", nil
	}

	f2 := func(
		f reflect.Kind,
		t reflect.Kind,
		data interface{}) (interface{}, error) {
		return data.(string) + "cat", nil
	}

	f := ComposeDecodeHookFunc(f1, f2)

	result, err := DecodeHookExec(
		f, reflect.TypeOf(""), reflect.TypeOf([]byte("")), "")
	if err != nil {
		t.Fatalf("bad: %s", err)
	}
	if result.(string) != "foobar" {
		t.Fatalf("bad: %#v", result)
	}
}
