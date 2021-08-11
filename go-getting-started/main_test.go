package main

import "testing"

func TestJokes(t *testing.T){
		var tests = []struct {
			input    int
			expected int
		}{
			{100000000000, 1},
			{100000000000, 100000000000},
			{"goodbye"},
			{ , }
		
		}
	
		for _, test := range tests {
			if output := getJokes(test.input); output != test.expected {
				t.Error("Test Failed: {} inputted, {} expected, recieved: {}", test.input, test.expected, output)
			}
		}
}