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
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("error, not expected word lenght: %d != %d", len(actual), len(c.expected))
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
