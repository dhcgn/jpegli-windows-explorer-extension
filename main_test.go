package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dhcgn/jpegli-windows-explorer-extension/settings"
)

func defaultTestingSettings() *settings.Seetings {
	return &settings.Seetings{
		// Default values
		Distance:             0.5,
		OverrideOriginalFile: false,

		// Because this is for testing, set these to true
		SkipUpdateCheck:   true,
		NoUserInteraction: true,
	}
}

func TestRun_Help(t *testing.T) {
	args := []string{"app", "--help"}

	exitCode := Run(args, defaultTestingSettings())

	if exitCode != ExitCodeSuccess {
		t.Errorf("Expected exit code %d, got %d", ExitCodeSuccess, exitCode)
	}
}

func prepareTestFile(t *testing.T, originalFile string) (string, func()) {
	content, err := os.ReadFile(originalFile)
	if err != nil {
		t.Fatalf("Failed to read original file: %v", err)
	}

	tempFile, err := os.CreateTemp("test-files", "test-image-*.jpg")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tempFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tempFile.Name(), func() {
		os.Remove(tempFile.Name())
	}
}

func TestRun_ConvertFile(t *testing.T) {
	originalFile := "test-files\\DSC_4045-NEF_DxO_DeepPRIME.jpg"
	testFile, cleanup := prepareTestFile(t, originalFile)
	defer cleanup()

	// Calculate expected output filename
	baseName := filepath.Base(testFile)
	ext := filepath.Ext(baseName)
	targetName := strings.TrimSuffix(baseName, ext) + ".jpegli.jpg"
	expectedOutputFile := filepath.Join(filepath.Dir(testFile), targetName)

	// Cleanup output before test (just in case) and after
	os.Remove(expectedOutputFile)
	defer os.Remove(expectedOutputFile)

	args := []string{"app", testFile}

	exitCode := Run(args, defaultTestingSettings())

	if exitCode != ExitCodeSuccess {
		t.Errorf("Expected exit code %d, got %d", ExitCodeSuccess, exitCode)
	}

	if _, err := os.Stat(expectedOutputFile); os.IsNotExist(err) {
		t.Errorf("Expected output file %s to be created", expectedOutputFile)
	} else {
		// Check if new file is smaller
		inputStat, err := os.Stat(testFile)
		if err != nil {
			t.Fatalf("Failed to stat input file: %v", err)
		}
		outputStat, err := os.Stat(expectedOutputFile)
		if err != nil {
			t.Fatalf("Failed to stat output file: %v", err)
		}

		if outputStat.Size() >= inputStat.Size() {
			t.Errorf("Expected output file to be smaller than input file. Input: %d, Output: %d", inputStat.Size(), outputStat.Size())
		}
	}
}

func TestRun_ConvertFile_Override(t *testing.T) {
	originalFile := "test-files\\DSC_4045-NEF_DxO_DeepPRIME.jpg"
	testFile, cleanup := prepareTestFile(t, originalFile)
	defer cleanup()

	// Get original size
	inputStat, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat input file: %v", err)
	}
	originalSize := inputStat.Size()

	// Calculate potential side-effect filename (which should NOT exist)
	baseName := filepath.Base(testFile)
	ext := filepath.Ext(baseName)
	targetName := strings.TrimSuffix(baseName, ext) + ".jpegli.jpg"
	unexpectedOutputFile := filepath.Join(filepath.Dir(testFile), targetName)

	// Ensure unexpected file doesn't exist before run
	os.Remove(unexpectedOutputFile)
	defer os.Remove(unexpectedOutputFile)

	args := []string{"app", testFile}

	opts := defaultTestingSettings()
	opts.OverrideOriginalFile = true

	exitCode := Run(args, opts)

	if exitCode != ExitCodeSuccess {
		t.Errorf("Expected exit code %d, got %d", ExitCodeSuccess, exitCode)
	}

	// Check if original file still exists (it should be overwritten)
	outputStat, err := os.Stat(testFile)
	if os.IsNotExist(err) {
		t.Errorf("Expected original file %s to exist (overwritten)", testFile)
	} else if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	} else {
		// Check if new file is smaller
		if outputStat.Size() >= originalSize {
			t.Errorf("Expected overwritten file to be smaller. Original: %d, New: %d", originalSize, outputStat.Size())
		}
	}

	// Check that no .jpegli.jpg file was created
	if _, err := os.Stat(unexpectedOutputFile); !os.IsNotExist(err) {
		t.Errorf("Expected no .jpegli.jpg file to be created, but found %s", unexpectedOutputFile)
	}
}
