package utils

import (
	"testing"
)

func TestGetMachineID(t *testing.T) {
	machineid := GetMachineID()
	if machineid == "" {
		t.Error("GetMachineID() failed")
	}
}

func TestGetMachineIDUseSHA256(t *testing.T) {
	machineid := GetMachineIDUseSHA256()
	if machineid == "" {
		t.Error("GetMachineIDUseSHA256() failed")
	}
}
