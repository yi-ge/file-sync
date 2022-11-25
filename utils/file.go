package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// CurrentFile get the current file path
func CurrentFile() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return file
}

// CurrentDir get the directory of the current executable file
func CurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// MakeDirIfNotExist create if the folder does not exist
func MakeDirIfNotExist(path string) {
	isExist, err := FileExists(path)
	if isExist && err != nil {
		os.MkdirAll(path, 0755)
	}
}

// FileExists checks if a file exists and returns a boolean or an error
func FileExists(fileName string) (bool, error) {
	if _, err := os.Stat(fileName); err == nil {
		// path/to/whatever exists
		return true, nil
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		return false, nil
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false, err
	}
}

// ReadFileToBytes will return the contents of a file as a byte slice
func ReadFileToBytes(path string) ([]byte, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(absolutePath)
}

// WriteByteFile creates a file from a byte slice with an optional file mode, only if it's new, and populates it - can force overwrite optionally
func WriteByteFile(path string, content []byte, mode int, overwrite bool) (bool, error) {
	var fileMode os.FileMode
	if mode == 0 {
		fileMode = os.FileMode(0600)
	} else {
		fileMode = os.FileMode(mode)
	}

	fileCheck, err := FileExists(path)
	check(err)

	// If not, create one with a starting digit
	if !fileCheck {
		err = os.WriteFile(path, content, fileMode)
		check(err)
		return true, err
	}

	// If the file exists and we want to overwrite it
	if fileCheck && overwrite {
		err = os.WriteFile(path, content, fileMode)
		check(err)
		return true, err
	}
	return false, nil
}
