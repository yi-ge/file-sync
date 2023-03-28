package utils

import (
	"bytes"
	"os"
	"testing"
)

func TestWriteRSAKeyPair(t *testing.T) {
	// create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// generate key pair buffers for testing
	privKeyBuf := bytes.NewBufferString("private key buffer")
	pubKeyBuf := bytes.NewBufferString("public key buffer")

	// call function being tested
	MakeDirIfNotExist(tmpDir + "/keypair/")
	privKeyFile, pubKeyFile, err := WriteRSAKeyPair(privKeyBuf, pubKeyBuf, tmpDir+"/keypair/")

	// check if any errors occurred during execution
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// check if files were created successfully
	if privKeyFile == false || pubKeyFile == false {
		t.Errorf("one or both of the key files failed to be created")
	}

	// check if files have correct permissions
	privKeyFileInfo, err := os.Stat(tmpDir + "/keypair/.priv.pem")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	pubKeyFileInfo, err := os.Stat(tmpDir + "/keypair/.pub.pem")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if privKeyFileInfo.Mode().Perm() != 0400 || pubKeyFileInfo.Mode().Perm() != 0644 {
		t.Errorf("one or both of the key files have incorrect permissions")
	}
}

func TestDeleteRSAKeyPair(t *testing.T) {
	// create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// create temporary key pair files for testing
	// generate key pair buffers for testing
	privKeyBuf := bytes.NewBufferString("private key buffer")
	pubKeyBuf := bytes.NewBufferString("public key buffer")

	// call function being tested
	MakeDirIfNotExist(tmpDir + "/keypair/")
	_, _, err = WriteRSAKeyPair(privKeyBuf, pubKeyBuf, tmpDir+"/keypair/")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// call function being tested
	err = DeleteRSAKeyPair(tmpDir + "/keypair/")

	// check if any errors occurred during execution
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// check if files were deleted successfully
	if _, err := os.Stat(tmpDir + "/keypair/.priv.pem"); !os.IsNotExist(err) {
		t.Errorf("private key file was not deleted")
	}
	if _, err := os.Stat(tmpDir + "/keypair/.pub.pem"); !os.IsNotExist(err) {
		t.Errorf("public key file was not deleted")
	}
}
