// Package main provides ...
package main

import (
	"encoding/json"
	"testing"
)

func TestJsonUnmarshal(t *testing.T) {
	var b = []byte(`{"a": "b"}`)
	var data interface{}
	err := json.Unmarshal(b, &data)
	t.Error(data)
	if err != nil {
		t.Error(err)
	}
}
