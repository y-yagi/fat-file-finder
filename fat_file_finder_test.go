package main

import (
	"bytes"
	"testing"

	"code.cloudfoundry.org/bytefmt"
)

func TestSearch(t *testing.T) {
	thresholdSize, _ := bytefmt.ToBytes("100K")
	out := new(bytes.Buffer)
	search("testdata", int64(thresholdSize), out)

	expected := "testdata/1m (1000K)\n"
	if out.String() != expected {
		t.Fatalf("expected '%s' but '%s'", expected, out.String())
	}
}
