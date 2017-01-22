package main

import (
	"testing"

	"github.com/cloudfoundry/bytefmt"
)

func TestSearch(t *testing.T) {
	thresholdSize, _ := bytefmt.ToBytes("100K")
	chann := search("testdata", int64(thresholdSize))

	expected := "testdata/1m (1000K)"
	for msg := range chann {
		if msg != expected {
			t.Fatalf("expected %s but %s", expected, msg)
		}
	}
}
