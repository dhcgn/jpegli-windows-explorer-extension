package main

import "testing"

func TestShouldSkipAlreadyProcessed(t *testing.T) {
	tests := []struct {
		name        string
		optimizedBy string
		want        bool
	}{
		{
			name:        "skip when marker exists",
			optimizedBy: "jpegli-windows-explorer-extension 1.0.0",
			want:        true,
		},
		{
			name:        "do not skip when marker is empty",
			optimizedBy: "",
			want:        false,
		},
		{
			name:        "do not skip when marker is whitespace only",
			optimizedBy: "   ",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSkipAlreadyProcessed(tt.optimizedBy)
			if got != tt.want {
				t.Fatalf("shouldSkipAlreadyProcessed(%q) = %v, want %v", tt.optimizedBy, got, tt.want)
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
