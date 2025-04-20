//go:build integration
// +build integration

package install

import (
	_ "embed"
	"testing"
)

func TestExtractEmbeddedZipFilesToAppFolder(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "TestExtractEmbeddedZipFilesToAppFolder"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExtractEmbeddedZipFilesToAppFolder()
		})
	}
}
