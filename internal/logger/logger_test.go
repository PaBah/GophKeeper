package logger

import (
	"testing"
)

func TestLog(t *testing.T) {
	cases := []struct {
		testName        string
		initialize      string
		isErrorExpected bool
	}{
		{"Test default level", "", false},
		{"Test debug level", "debug", false},
		{"Test info level", "info", false},
		{"Test warn level", "warn", false},
		{"Test error level", "error", false},
		{"Test invalid level", "invalid", true},
		{"Test capitalized level", "DEBUG", false},
		{"Test mixed case level", "iNfO", false},
	}

	for _, tt := range cases {
		t.Run(tt.testName, func(t *testing.T) {
			err := Initialize(tt.initialize)
			got := Log()

			if tt.isErrorExpected && err == nil {
				t.Errorf("Got no error, but expected one")
			}

			if !tt.isErrorExpected && err != nil {
				t.Errorf("Got error %v, but not expected one", err)
			}

			if !tt.isErrorExpected && err != nil {
				t.Errorf("Got error %v, but not expected one", err)
			}

			if !tt.isErrorExpected && got == nil {
				t.Errorf("Got nil, wanted non-nil logger instance")
			}

		})
	}
}
