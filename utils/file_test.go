package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCurrentFile(t *testing.T) {
	filename := CurrentFile()
	expectedFilename := filepath.Join("utils", "utils_test.go")
	if strings.Contains(expectedFilename, filename) {
		t.Errorf("CurrentFile() returned %q, expected %q", filename, expectedFilename)
	}
}

func TestCurrentDir(t *testing.T) {
	dirname := CurrentDir()
	expectedDirname, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}
	if dirname != expectedDirname {
		t.Errorf("CurrentDir() returned %q, expected %q", dirname, expectedDirname)
	}
}

func TestMakeDirIfNotExist(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	newDir := filepath.Join(tempDir, "newdir")
	MakeDirIfNotExist(newDir)
	if _, err := os.Stat(newDir); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%q should exist after calling MakeDirIfNotExist(%q)", newDir, newDir)
		} else {
			t.Errorf("Error checking if %q exists after calling MakeDirIfNotExist(%q): %v", newDir, newDir, err)
		}
	}

	MakeDirIfNotExist(newDir)
	if _, err := os.Stat(newDir); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("%q should exist after calling MakeDirIfNotExist(%q) a second time", newDir, newDir)
		} else {
			t.Errorf("Error checking if %q exists after calling MakeDirIfNotExist(%q) a second time: %v", newDir, newDir, err)
		}
	}
}

func TestFileExists(t *testing.T) {
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if exists, err := FileExists(tempFile.Name()); !exists || err != nil {
		t.Errorf("%q should exist, but FileExists(%q) returned (%t, %v)", tempFile.Name(), tempFile.Name(), exists, err)
	}

	if exists, err := FileExists(filepath.Join(tempFile.Name(), "nonexistent")); exists || err == nil {
		t.Errorf("%q should not exist, but FileExists(%q) returned (%t, %v)", filepath.Join(tempFile.Name(), "nonexistent"), filepath.Join(tempFile.Name(), "nonexistent"), exists, err)
	}
}

func TestReadFileToBytes(t *testing.T) {
	expectedContent := []byte("test content")
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	err = os.WriteFile(tempFile.Name(), expectedContent, 0666)
	if err != nil {
		t.Fatalf("Error writing to temp file: %v", err)
	}

	content, err := ReadFileToBytes(tempFile.Name())
	if err != nil {
		t.Errorf("Error reading file %q: %v", tempFile.Name(), err)
	}
	if !bytes.Equal(content, expectedContent) {
		t.Errorf("Content of file %q is %q, expected %q", tempFile.Name(), content, expectedContent)
	}
}

func TestWriteByteFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "testfile")
	content := []byte("test content")

	// Write to new file with default mode
	written, err := WriteByteFile(tempFile, content, 0, false)
	if err != nil {
		t.Errorf("Error writing to file %q: %v", tempFile, err)
	}
	if !written {
		t.Errorf("WriteByteFile(%q, _, _, _) returned false, but should have created a new file", tempFile)
	}
	fileInfo, err := os.Stat(tempFile)
	if err != nil {
		t.Fatalf("Error getting file info for %q: %v", tempFile, err)
	}
	expectedMode := os.FileMode(0600)
	if fileInfo.Mode() != expectedMode {
		t.Errorf("Newly创建的文件 %q 的模式应该是 %v，但实际上是 %v", tempFile, expectedMode, fileInfo.Mode())
	}

	// Try to write to existing file without overwriting
	written, err = WriteByteFile(tempFile, content, 0, false)
	if err != nil {
		t.Errorf("Error writing to file %q: %v", tempFile, err)
	}
	if written {
		t.Errorf("WriteByteFile(%q, _, _, false) returned true, but should not have overwritten the existing file", tempFile)
	}

	// Try to write to existing file with overwriting
	written, err = WriteByteFile(tempFile, content, 0, true)
	if err != nil {
		t.Errorf("Error writing to file %q: %v", tempFile, err)
	}
	if !written {
		t.Errorf("WriteByteFile(%q, _, _, true) returned false, but should have overwritten the existing file", tempFile)
	}
	newContent, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Error reading file %q after overwriting it: %v", tempFile, err)
	}
	if !bytes.Equal(newContent, content) {
		t.Errorf("Content of file %q after overwriting is %q, expected %q", tempFile, newContent, content)
	}
}

func TestFileSHA256(t *testing.T) {
	expectedHash := "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	err = os.WriteFile(tempFile.Name(), []byte("test content"), 0666)
	if err != nil {
		t.Fatalf("Error writing to temp file: %v", err)
	}
	hash, err := FileSHA256(tempFile.Name())
	if err != nil {
		t.Errorf("Error computing SHA256 hash of file %q: %v", tempFile.Name(), err)
	}
	if hash != expectedHash {
		t.Errorf("SHA256 hash of file %q is %q, expected %q", tempFile.Name(), hash, expectedHash)
	}

	_, err = FileSHA256(filepath.Join(tempFile.Name(), "nonexistent"))
	if err == nil {
		t.Error("FileSHA256 should have returned an error for a nonexistent file")
	}
}
