package alex_test

import (
	"testing"

	"example.com/go-demo-1/mascot"
)

func TestMascot(t *testing.T) {
	if mascot.BestAlex() != "Alex_Best" {
		t.Fatal("Wrong test :(")
	}
}
