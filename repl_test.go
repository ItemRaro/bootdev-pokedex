package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Lorem Ipsum is simply dummy text",
			expected: []string{"lorem", "ipsum", "is", "simply", "dummy", "text"},
		},
		{
			input:    "   Lorem ipsum dolor sit amet",
			expected: []string{"lorem", "ipsum", "dolor", "sit", "amet"},
		},
		{
			input:    "consectetur adipiscing elit   ",
			expected: []string{"consectetur", "adipiscing", "elit"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "Catch\tPIKACHU\n",
			expected: []string{"catch", "pikachu"},
		},
		{
			input:    "  EXPLORE   canalave-city-area  ",
			expected: []string{"explore", "canalave-city-area"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("error, not expected word lenght: %d != %d", len(actual), len(c.expected))
			// Comparing words below would index past the end of the shorter slice.
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("error, not the expected word: %s != %s", word, expectedWord)
			}
		}
	}
}
