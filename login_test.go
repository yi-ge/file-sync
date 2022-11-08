package main

import (
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Log(os.Environ())
	login()
}

func TestGetMachineID(t *testing.T) {
	t.Log(getMachineID())
}
