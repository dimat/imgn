package models

import "testing"

func TestResolve(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantID  string
		wantErr bool
	}{
		{"flash2 alias", "flash2", Flash2.ID, false},
		{"pro alias", "pro", Pro.ID, false},
		{"flash alias", "flash", Flash.ID, false},
		{"flash2 full ID", "gemini-3.1-flash-image-preview", Flash2.ID, false},
		{"pro full ID", "gemini-3-pro-image-preview", Pro.ID, false},
		{"flash full ID", "gemini-2.0-flash-preview-image-generation", Flash.ID, false},
		{"unknown", "gpt-4", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Resolve(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m.ID != tt.wantID {
				t.Errorf("got %s, want %s", m.ID, tt.wantID)
			}
		})
	}
}

func TestValidateAspect(t *testing.T) {
	// Common aspects valid for all models
	common := []string{"1:1", "16:9", "9:16", "4:3", "3:4"}
	for _, a := range common {
		if err := ValidateAspect(a, Pro); err != nil {
			t.Errorf("expected %s to be valid for Pro: %v", a, err)
		}
		if err := ValidateAspect(a, Flash2); err != nil {
			t.Errorf("expected %s to be valid for Flash2: %v", a, err)
		}
	}
	// Flash2-only aspects
	flash2Only := []string{"1:4", "4:1", "1:8", "8:1"}
	for _, a := range flash2Only {
		if err := ValidateAspect(a, Flash2); err != nil {
			t.Errorf("expected %s to be valid for Flash2: %v", a, err)
		}
		if err := ValidateAspect(a, Pro); err == nil {
			t.Errorf("expected %s to be invalid for Pro", a)
		}
	}
	if err := ValidateAspect("2:1", Pro); err == nil {
		t.Error("expected 2:1 to be invalid")
	}
}

func TestValidateSize(t *testing.T) {
	if err := ValidateSize("4k", Pro); err != nil {
		t.Errorf("4k should be valid for Pro: %v", err)
	}
	if err := ValidateSize("4k", Flash); err == nil {
		t.Error("4k should be invalid for Flash")
	}
	if err := ValidateSize("2k", Flash); err != nil {
		t.Errorf("2k should be valid for Flash: %v", err)
	}
	// Flash2 supports 512px and 4k
	if err := ValidateSize("512px", Flash2); err != nil {
		t.Errorf("512px should be valid for Flash2: %v", err)
	}
	if err := ValidateSize("4k", Flash2); err != nil {
		t.Errorf("4k should be valid for Flash2: %v", err)
	}
	if err := ValidateSize("512px", Pro); err == nil {
		t.Error("512px should be invalid for Pro")
	}
}
