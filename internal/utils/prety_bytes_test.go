package utils

import "testing"

func TestHumanReadableSize(t *testing.T) {
	testCases := []struct {
		name           string
		input          uint64
		expectedOutput string
	}{
		{
			name:           "Bytes",
			input:          500,
			expectedOutput: "500 B",
		},
		{
			name:           "KiloBytes",
			input:          1024,
			expectedOutput: "1.0 KB",
		},
		{
			name:           "MegaBytes",
			input:          1048576, // 1024^2
			expectedOutput: "1.0 MB",
		},
		{
			name:           "GigaBytes",
			input:          1073741824, // 1024^3
			expectedOutput: "1.0 GB",
		},
		{
			name:           "TeraBytes",
			input:          1099511627776, // 1024^4
			expectedOutput: "1.0 TB",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualOutput := HumanReadableSize(testCase.input)
			if actualOutput != testCase.expectedOutput {
				t.Errorf("Test %s failed: input %d, expected output %s, but got %s",
					testCase.name, testCase.input, testCase.expectedOutput, actualOutput)
			}
		})
	}
}
