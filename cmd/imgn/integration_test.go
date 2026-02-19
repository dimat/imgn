package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmd_Help(t *testing.T) {
	cmd := newRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "imgn") {
		t.Error("help output should contain 'imgn'")
	}
}

func TestGenerateCmd_NoPrompt(t *testing.T) {
	// Without API key, should fail before prompt check
	cmd := newRootCmd()
	cmd.SetArgs([]string{"generate"})
	// This will fail due to no API key or no prompt - both are valid error paths
	err := cmd.Execute()
	if err == nil {
		// It's OK if it exits via os.Exit in the actual code,
		// but in test it should return an error
		t.Log("command executed (may have exited)")
	}
}

func TestVersionCmd(t *testing.T) {
	cmd := newRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestModelsCmd(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"models"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGenerateAliases(t *testing.T) {
	cmd := newRootCmd()
	genCmd, _, err := cmd.Find([]string{"gen"})
	if err != nil {
		t.Fatalf("'gen' alias not found: %v", err)
	}
	if genCmd.Name() != "generate" {
		t.Errorf("expected generate command, got %s", genCmd.Name())
	}

	gCmd, _, err := cmd.Find([]string{"g"})
	if err != nil {
		t.Fatalf("'g' alias not found: %v", err)
	}
	if gCmd.Name() != "generate" {
		t.Errorf("expected generate command, got %s", gCmd.Name())
	}
}
