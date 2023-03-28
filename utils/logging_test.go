package utils

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	defer log.SetOutput(log.Writer())

	// call function being tested with a non-nil error
	// define and initialize someError
	someError := errors.New("someError")
	Check(someError)
	loggedMsg := buf.String()
	expectedMsg := "someError"
	if !strings.Contains(loggedMsg, expectedMsg) {
		t.Errorf("unexpected logged message: %q", loggedMsg)
	}

	// call function being tested with a nil error
	buf.Reset()
	Check(nil)
	loggedMsg = buf.String()
	if loggedMsg != "" {
		t.Errorf("unexpected logged message: %q", loggedMsg)
	}
}
