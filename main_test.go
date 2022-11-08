package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	t.Log(os.Environ())
}
