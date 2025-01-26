package sifo

import (
	"testing"
)

func TestIsCloseMatch(t *testing.T) {
	words := map[string]int64{
		"apple":  1,
		"banana": 1,
		"cherry": 1,
		"date":   1,
		"fig":    1,
		"grape":  1,
	}

	tests := []struct {
		word     string
		expected bool
	}{
		{"apple", true},   // Exact match
		{"appl", true},    // One letter off
		{"apples", true},  // One letter off
		{"applf", true},   // One letter off
		{"applz", true},   // One letter off
		{"banan", true},   // One letter off
		{"bananas", true}, // One letter off
		{"bananz", true},  // One letter off
		{"grap", true},    // One letter off
		{"grapz", true},   // One letter off
		{"bonona", true},  // Two letters off, 6 letters
		{"trope", false},  // Two letters off
		{"xyz", false},    // No match
	}

	for _, test := range tests {
		result := isCloseMatch(test.word, words)
		if result != test.expected {
			t.Errorf("isCloseMatch(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}

func TestHasVowel(t *testing.T) {
	tests := []struct {
		word     string
		expected bool
	}{
		{"apple", true},   // Contains vowels 'a' and 'e'
		{"sky", false},    // No vowels
		{"banana", true},  // Contains vowels 'a'
		{"rhythm", false}, // No vowels
		{"grape", true},   // Contains vowels 'a' and 'e'
		{"fly", false},    // No vowels
		{"queue", true},   // Contains vowels 'u' and 'e'
		{"", false},       // Empty string
		{"bcdfg", false},  // No vowels
		{"aeiou", true},   // Contains all vowels
	}

	for _, test := range tests {
		result := hasVowel(test.word)
		if result != test.expected {
			t.Errorf("hasVowel(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}

func TestHasPrefix(t *testing.T) {
	prefixes := map[string]bool{
		"ap":  true,
		"ban": true,
		"che": true,
		"da":  true,
		"fi":  true,
		"gr":  true,
	}

	tests := []struct {
		word     string
		expected bool
	}{
		{"apple", true},      // Matches 2-letter prefix "ap"
		{"banana", true},     // Matches 3-letter prefix "ban"
		{"cherry", true},     // Matches 3-letter prefix "che"
		{"date", true},       // Matches 2-letter prefix "da"
		{"fig", true},        // Matches 2-letter prefix "fi"
		{"grape", true},      // Matches 2-letter prefix "gr"
		{"grapefruit", true}, // Matches 2-letter prefix "gr"
		{"apricot", true},    // Matches 2-letter prefix "ap"
		{"berry", false},     // No matching prefix
		{"", false},          // Empty string
		{"a", false},         // Single letter
		{"an", false},        // Two letters, no matching prefix
		{"ban", false},       // Matches 3-letter prefix "ban" but only 3 letters long
		{"bana", true},       // Matches 3-letter prefix "ban"
		{"b", false},         // Single letter, no matching prefix
	}

	for _, test := range tests {
		result := hasPrefix(test.word, prefixes)
		if result != test.expected {
			t.Errorf("hasPrefix(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}

func TestHasSuffix(t *testing.T) {
	suffixes := map[string]bool{
		"le":  true,
		"ana": true,
		"rry": true,
		"te":  true,
		"ig":  true,
		"pe":  true,
	}

	tests := []struct {
		word     string
		expected bool
	}{
		{"apple", true},       // Matches 2-letter suffix "le"
		{"banana", true},      // Matches 3-letter suffix "ana"
		{"cherry", true},      // Matches 3-letter suffix "rry"
		{"date", true},        // Matches 2-letter suffix "te"
		{"fig", true},         // Matches 2-letter suffix "ig"
		{"grape", true},       // Matches 2-letter suffix "pe"
		{"grapefruit", false}, // No matching suffix
		{"apricot", false},    // No matching suffix
		{"berry", true},       // Matches 3-letter suffix "rry"
		{"", false},           // Empty string
		{"a", false},          // Single letter
		{"an", false},         // Two letters, no matching suffix
		{"ana", false},        // Matches 3-letter suffix "ana" but only 3 letters long
		{"anana", true},       // Matches 3-letter suffix "ana"
		{"b", false},          // Single letter, no matching suffix
	}

	for _, test := range tests {
		result := hasSuffix(test.word, suffixes)
		if result != test.expected {
			t.Errorf("hasSuffix(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}

func TestHasMiddles(t *testing.T) {
	middles := map[string]bool{
		"pp":  true,
		"pl":  true,
		"ppl": true,
	}

	tests := []struct {
		word     string
		expected bool
	}{
		{"apple", true},   // All middles ("pp", "pl", "ppl") match
		{"apples", false}, // Not all middles match
		{"", false},       // Empty string
	}

	for _, test := range tests {
		result := hasMiddles(test.word, middles)
		if result != test.expected {
			t.Errorf("hasMiddles(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}
