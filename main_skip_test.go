package main

import "testing"

func TestShouldSkipAlreadyProcessed(t *testing.T) {
	tests := []struct {
		name                 string
		alwaysReprocessFiles bool
		optimizedBy          string
		want                 bool
	}{
		{
			name:                 "skip when marker exists and always reprocess disabled",
			alwaysReprocessFiles: false,
			optimizedBy:          "jpegli-windows-explorer-extension 1.0.0",
			want:                 true,
		},
		{
			name:                 "do not skip when marker is empty",
			alwaysReprocessFiles: false,
			optimizedBy:          "",
			want:                 false,
		},
		{
			name:                 "do not skip when marker is whitespace only",
			alwaysReprocessFiles: false,
			optimizedBy:          "   ",
			want:                 false,
		},
		{
			name:                 "do not skip when always reprocess enabled",
			alwaysReprocessFiles: true,
			optimizedBy:          "jpegli-windows-explorer-extension 1.0.0",
			want:                 false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSkipAlreadyProcessed(tt.alwaysReprocessFiles, tt.optimizedBy)
			if got != tt.want {
				t.Fatalf("shouldSkipAlreadyProcessed(%v, %q) = %v, want %v", tt.alwaysReprocessFiles, tt.optimizedBy, got, tt.want)
			}
		})
	}
}

func TestOptimizedByValueIncludesVersion(t *testing.T) {
	previousVersion := Version
	Version = "9.9.9"
	t.Cleanup(func() { Version = previousVersion })

	got := optimizedByValue()
	want := "jpegli-windows-explorer-extension 9.9.9"
	if got != want {
		t.Fatalf("optimizedByValue() = %q, want %q", got, want)
	}
}
