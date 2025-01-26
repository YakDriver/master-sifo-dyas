package sifo

import (
	"reflect"
	"testing"
)

func TestWordPattern(t *testing.T) {
	tests := []struct {
		word     string
		expected string
	}{
		{"at", "vc"},         // Vowel followed by consonant
		{"street", "cvc"},    // Consonant group, vowel group, consonant group
		{"apple", "vcv"},     // Vowel group, consonant group, vowel group
		{"banana", "cvcvcv"}, // Vowel group, consonant group, vowel group, consonant group, vowel group
		{"rhythm", "cvc"},    // Consonant group, vowel group, consonant group
		{"", ""},             // Empty string
		{"a", "v"},           // Single vowel
		{"b", "c"},           // Single consonant
		{"aeiou", "v"},       // Single vowel group
		{"bcdfg", "c"},       // Single consonant group
		{"aeioubcdfg", "vc"}, // Vowel group followed by consonant group
		{"city", "cvcv"},     // Vowel group, consonant group, vowel group
		{"sky", "cv"},        // Consonant group, vowel group, consonant group
		{"yellow", "vcvc"},   // Vowel group, consonant group, vowel group
	}

	for _, test := range tests {
		result := wordPattern(test.word)
		if result != test.expected {
			t.Errorf("wordPattern(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}

func TestVowelGroups(t *testing.T) {
	tests := []struct {
		word     string
		expected []string
	}{
		{"apple", []string{"a", "e"}},            // Single vowels separated by consonants
		{"banana", []string{"a"}},                // Single vowels separated by consonants
		{"beautiful", []string{"eau", "i", "u"}}, // Multiple vowels in groups
		{"rhythm", []string{"y"}},                // Single vowel group
		{"", nil},                                // Empty string
		{"a", []string{"a"}},                     // Single vowel
		{"b", nil},                               // Single consonant
		{"aeiou", []string{"aeiou"}},             // Single vowel group
		{"bcdfg", nil},                           // No vowels
		{"bcdfgaeioubcdfg", []string{"aeiou"}},   // Vowel group followed by consonants
		{"milwaukee", []string{"i", "au", "ee"}}, // Vowel group followed by consonants
	}

	for _, test := range tests {
		result := vowelGroups(test.word)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("vowelGroups(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}

func TestConsonantGroups(t *testing.T) {
	tests := []struct {
		word     string
		expected []string
	}{
		{"apple", []string{"ppl"}},                  // Consonant group
		{"banana", []string{"b", "n"}},              // Single consonants separated by vowels
		{"beautiful", []string{"b", "t", "f", "l"}}, // Single consonants separated by vowels
		{"rhythm", []string{"rh", "thm"}},           // Consonant groups
		{"", nil},                                   // Empty string
		{"a", nil},                                  // Single vowel
		{"b", []string{"b"}},                        // Single consonant
		{"aeiou", nil},                              // No consonants
		{"bcdfg", []string{"bcdfg"}},                // Single consonant group
		{"aeioubcdfgaeiou", []string{"bcdfg"}},      // Vowels followed by consonant group
		{"milwaukee", []string{"m", "lw", "k"}},     // Consonant groups separated by vowels
	}

	for _, test := range tests {
		result := consonantGroups(test.word)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("consonantGroups(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}

func TestVowelConsonantBoundaries(t *testing.T) {
	tests := []struct {
		word     string
		expected []string
	}{
		{"apple", []string{"ap"}},                 // Vowel followed by consonant
		{"street", []string{"et"}},                // Consonant followed by vowel
		{"banana", []string{"an"}},                // Vowel followed by consonant
		{"beautiful", []string{"ut", "if", "ul"}}, // Multiple boundaries
		{"rhythm", []string{"yt"}},                // Consonant followed by vowel
		{"", nil},                                 // Empty string
		{"a", nil},                                // Single vowel
		{"b", nil},                                // Single consonant
		{"aeiou", nil},                            // All vowels
		{"bcdfg", nil},                            // All consonants
		{"aeioubcdfg", []string{"ub"}},            // Vowel followed by consonant
	}

	for _, test := range tests {
		result := vowelConsonantBoundaries(test.word)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("vowelConsonantBoundaries(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}

func TestConsonantVowelBoundaries(t *testing.T) {
	tests := []struct {
		word     string
		expected []string
	}{
		{"apple", []string{"le"}},                 // Consonant followed by vowel
		{"street", []string{"re"}},                // Consonant followed by vowel
		{"banana", []string{"ba", "na"}},          // Consonant followed by vowel
		{"beautiful", []string{"be", "ti", "fu"}}, // Multiple boundaries
		{"rhythm", []string{"hy"}},                // Consonant followed by vowel
		{"", nil},                                 // Empty string
		{"a", nil},                                // Single vowel
		{"b", nil},                                // Single consonant
		{"aeiou", nil},                            // All vowels
		{"bcdfg", nil},                            // All consonants
		{"aeioubcdfg", nil},                       // No consonant followed by vowel
		{"bcdfgaeiou", []string{"ga"}},            // Consonant followed by vowel
	}

	for _, test := range tests {
		result := consonantVowelBoundaries(test.word)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("consonantVowelBoundaries(%q) = %v; want %v", test.word, result, test.expected)
		}
	}
}
