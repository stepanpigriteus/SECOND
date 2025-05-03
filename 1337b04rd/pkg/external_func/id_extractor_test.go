package externalfunc_test

import (
	"testing"

	externalfunc "a1337b04rd/pkg/external_func"
)

func TestExtractIDFromPath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantID    int
		expectErr bool
	}{
		{
			name:      "valid path",
			path:      "/api/character/42",
			wantID:    42,
			expectErr: false,
		},
		{
			name:      "invalid path - too short",
			path:      "/character",
			expectErr: true,
		},
		{
			name:      "invalid id - not a number",
			path:      "/api/character/abc",
			expectErr: true,
		},
		{
			name:      "trailing slash",
			path:      "/api/character/123/",
			expectErr: true, // т.к. последний сегмент пустой
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := externalfunc.ExtractIDFromPath(tt.path)
			if (err != nil) != tt.expectErr {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectErr && id != tt.wantID {
				t.Errorf("expected ID %d, got %d", tt.wantID, id)
			}
		})
	}
}
