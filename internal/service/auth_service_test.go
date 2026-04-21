package service_test

import (
	"testing"

	"github.com/bert1727/ChatApp/internal/service"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		plain   string
		want    string
		wantErr bool
	}{
		{"test1", "fasdf", "wasd", false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := service.HashPassword(tt.plain)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("HashPassword() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("HashPassword() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
