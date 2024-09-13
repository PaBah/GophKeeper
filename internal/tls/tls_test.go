package tls

import (
	"os"
	"testing"
)

func TestCreateTLSCert(t *testing.T) {
	tests := []struct {
		name        string
		certPath    string
		keyPath     string
		errExpected bool
	}{
		{
			name:        "ValidPaths",
			certPath:    "testCert.cert",
			keyPath:     "testKey.key",
			errExpected: false,
		},
		{
			name:        "InvalidCertPath",
			certPath:    "/non/existing/path/testCert.cert",
			keyPath:     "testKey.key",
			errExpected: true,
		},
		{
			name:        "InvalidKeyPath",
			certPath:    "testCert.cert",
			keyPath:     "/non/existing/path/testKey.key",
			errExpected: true,
		},
		{
			name:        "EmptyCertPath",
			certPath:    "",
			keyPath:     "testKey.key",
			errExpected: true,
		},
		{
			name:        "EmptyKeyPath",
			certPath:    "testCert.cert",
			keyPath:     "",
			errExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CreateTLSCert(test.certPath, test.keyPath)
			if (err != nil) != test.errExpected {
				t.Fatalf("Error: Expected error state %v, but got %v", test.errExpected, err != nil)
			}
		})
		// Cleanup after Test
		if _, err := os.Stat(test.certPath); err == nil {
			_ = os.Remove(test.certPath)
		}
		if _, err := os.Stat(test.keyPath); err == nil {
			_ = os.Remove(test.keyPath)
		}
	}
}
