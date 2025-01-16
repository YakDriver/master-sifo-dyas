package sifo

import (
	"testing"
)

func TestEncodeWord(t *testing.T) {
	cipher := Cipher{
		"a":   "x",
		"b":   "y",
		"c":   "z",
		"d":   "w",
		"ab":  "v",
		"bc":  "u",
		"abc": "n",
		"xy":  "ab",
		"yz":  "cd",
		"wz":  "ef",
		"e":   "abc",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"a", "x"},
		{"b", "y"},
		{"c", "z"},
		{"d", "w"},
		{"ab", "v"},
		{"bc", "u"},
		{"abc", "n"},
		{"aab", "xv"},
		{"abbc", "vu"},
		{"abca", "nx"},
		{"abcd", "nw"},
		{"Abcd", "Nw"},
		{"d", "w"},
		{"xy", "ab"},
		{"yz", "cd"},
		{"wz", "ef"},
		{"xyab", "abv"},
		{"XYAB", "ABV"},
		{"abxy", "vab"},
		{"xyyz", "abcd"},
		{"wzbc", "efu"},
		{"abcxy", "nab"},
		{"e", "abc"},
		{"eabc", "abcn"},
		{"eabca", "abcnx"},
	}

	for _, test := range tests {
		result := encodeWord(test.input, cipher)
		if result != test.expected {
			t.Errorf("encodeWord(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestDecodeWord(t *testing.T) {
	cipher := Cipher{
		"a":   "x",
		"b":   "y",
		"c":   "z",
		"d":   "w",
		"ab":  "v",
		"bc":  "u",
		"abc": "n",
		"xy":  "ab",
		"yz":  "cd",
		"wz":  "ef",
		"e":   "abc",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"x", "a"},
		{"y", "b"},
		{"z", "c"},
		{"w", "d"},
		{"v", "ab"},
		{"u", "bc"},
		{"n", "abc"},
		{"xv", "aab"},
		{"vu", "abbc"},
		{"nx", "abca"},
		{"nw", "abcd"},
		{"w", "d"},
		{"ab", "xy"},
		{"cd", "yz"},
		{"ef", "wz"},
		{"abv", "xyab"},
		{"vab", "abxy"},
		{"abcd", "xyyz"},
		{"efu", "wzbc"},
		{"nab", "abcxy"},
		{"abc", "e"},
		{"abcn", "eabc"},
		{"abcnx", "eabca"},
	}

	for _, test := range tests {
		result := decodeWord(test.input, cipher)
		if result != test.expected {
			t.Errorf("decodeWord(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestEncodeDecodeWordSimple(t *testing.T) {
	words := LoadWords("../words_with_rank.csv")

	cipher := kindCatsCipher()

	for word := range words {
		encodedWord := encodeWord(word, cipher)
		decodedWord := decodeWord(encodedWord, cipher)
		if decodedWord != word {
			t.Errorf("encodeWord and decodeWord mismatch: original %q, encoded %q, decoded %q", word, encodedWord, decodedWord)
		}
	}
}
