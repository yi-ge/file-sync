package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckURL(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		wantErr string
	}{
		{
			name:    "valid URL",
			urlStr:  "https://example.com/path",
			wantErr: "",
		},
		{
			name:    "missing scheme",
			urlStr:  "example.com/path",
			wantErr: "missing scheme",
		},
		{
			name:    "missing host",
			urlStr:  "https:///path",
			wantErr: "missing host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := CheckURL(tt.urlStr)
			if tt.wantErr == "" {
				assert.NoError(t, gotErr)
			} else {
				assert.EqualError(t, gotErr, tt.wantErr)
			}
		})
	}
}
