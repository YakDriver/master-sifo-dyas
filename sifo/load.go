package sifo

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadWords(filename string) map[string]int64 {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	words := make(map[string]int64)
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		i++
		if i == 1 {
			continue // skip first line
		}
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			continue // skip lines that don't have at least 2 parts
		}
		value, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue // skip lines where the second part is not a valid int64
		}
		words[parts[0]] = value
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d words\n", len(words))
	return words
}

func CreatePartialWordDictionary(words map[string]int64) map[string]int64 {
	partialWords := make(map[string]int64)

	for word, count := range words {
		length := len(word)
		for i := 0; i < length; i++ {
			for j := i + 1; j <= length; j++ {
				substr := word[i:j]
				partialWords[substr] += count
			}
		}
	}

	return partialWords
}

// ListWordPatterns breaks down words into vowel and consonant patterns. 1 or more vowels is a vowel group and
// 1 or more consonants is a consonant group. The function returns a map of patterns. For example, "at" would be
// "vc" because it has a vowel followed by a consonant. "street" would be "cvc" because it has a consonant group
// followed by a vowel group followed by a consonant group.
func WordPatterns(words map[string]int64) map[string]bool {
	patterns := make(map[string]bool)

	for word := range words {
		pattern := wordPattern(word)
		patterns[pattern] = true
	}

	return patterns
}

// wordPattern breaks down a word into a pattern of vowels and consonants. 1 or more vowels is a vowel group and
// 1 or more consonants is a consonant group. For example, "at" would be "vc" because it has a vowel followed by a
// consonant. "street" would be "cvc" because it has a consonant group followed by a vowel group followed by a
// consonant group.
func wordPattern(word string) string {
	vowels := "aeiouy"
	pattern := ""
	var lastCharType rune

	for _, char := range word {
		if strings.ContainsRune(vowels, char) {
			if lastCharType != 'v' {
				pattern += "v"
				lastCharType = 'v'
			}
		} else {
			if lastCharType != 'c' {
				pattern += "c"
				lastCharType = 'c'
			}
		}
	}

	return pattern
}

// ListVowelGroups returns a map of vowel groups found in the words. A vowel group is defined as 1 or more vowels.
// For now, this simply returns a map of all the vowel groups found in the words. This could be expanded to return
// counts of each vowel group.
func VowelGroups(words map[string]int64) map[string]bool {
	vowelGroupsMap := make(map[string]bool)

	for word := range words {
		groups := vowelGroups(word)
		for _, group := range groups {
			vowelGroupsMap[group] = true
		}
	}

	return vowelGroupsMap
}

// vowelGroups returns a slice of unique vowel groups found in the word. A vowel group is defined as 1 or more vowels.
func vowelGroups(word string) []string {
	vowels := "aeiouy"
	var groups []string
	var currentGroup string
	uniqueGroups := make(map[string]bool)

	for _, char := range word {
		if strings.ContainsRune(vowels, char) {
			currentGroup += string(char)
		} else {
			if currentGroup != "" {
				if !uniqueGroups[currentGroup] {
					groups = append(groups, currentGroup)
					uniqueGroups[currentGroup] = true
				}
				currentGroup = ""
			}
		}
	}

	if currentGroup != "" && !uniqueGroups[currentGroup] {
		groups = append(groups, currentGroup)
		uniqueGroups[currentGroup] = true
	}

	return groups
}

// ListConsonantGroups returns a map of consonant groups found in the words. A consonant group is defined as 1 or more consonants.
// For now, this simply returns a map of all the consonant groups found in the words. This could be expanded to return
// counts of each consonant group.
func ConsonantGroups(words map[string]int64) map[string]bool {
	consonantGroupsMap := make(map[string]bool)

	for word := range words {
		groups := consonantGroups(word)
		for _, group := range groups {
			consonantGroupsMap[group] = true
		}
	}

	return consonantGroupsMap
}

// consonantGroups returns a slice of unique consonant groups found in the word. A consonant group is defined as 1 or more consonants.
func consonantGroups(word string) []string {
	vowels := "aeiouy"
	var groups []string
	var currentGroup string
	uniqueGroups := make(map[string]bool)

	for _, char := range word {
		if !strings.ContainsRune(vowels, char) {
			currentGroup += string(char)
		} else {
			if currentGroup != "" {
				if !uniqueGroups[currentGroup] {
					groups = append(groups, currentGroup)
					uniqueGroups[currentGroup] = true
				}
				currentGroup = ""
			}
		}
	}

	if currentGroup != "" && !uniqueGroups[currentGroup] {
		groups = append(groups, currentGroup)
		uniqueGroups[currentGroup] = true
	}

	return groups
}

func VowelConsonantBoundaries(words map[string]int64) map[string]bool {
	boundaries := make(map[string]bool)

	for word := range words {
		boundary := vowelConsonantBoundaries(word)
		for _, b := range boundary {
			boundaries[b] = true
		}
	}

	return boundaries
}

// vowelConsonantBoundaries returns a slice of 2-length strings that represent the unique boundaries between vowel
// groups and consonant groups in the word. For example, "apple" would return ["ap"] because it is the only vowel
// group followed by a consonant group, "street" would return ["et"], "banana" would return ["an"], and
// "beautiful" would return ["ut", "if", "ul"].
func vowelConsonantBoundaries(word string) []string {
	vowels := "aeiouy"
	var boundaries []string
	uniqueBoundaries := make(map[string]bool)

	for i := 0; i < len(word)-1; i++ {
		currentChar := word[i]
		nextChar := word[i+1]

		currentIsVowel := strings.ContainsRune(vowels, rune(currentChar))
		nextIsVowel := strings.ContainsRune(vowels, rune(nextChar))

		if currentIsVowel && !nextIsVowel {
			boundary := string(currentChar) + string(nextChar)
			if !uniqueBoundaries[boundary] {
				boundaries = append(boundaries, boundary)
				uniqueBoundaries[boundary] = true
			}
		}
	}

	return boundaries
}

func ConsonantVowelBoundaries(words map[string]int64) map[string]bool {
	boundaries := make(map[string]bool)

	for word := range words {
		boundary := consonantVowelBoundaries(word)
		for _, b := range boundary {
			boundaries[b] = true
		}
	}

	return boundaries
}

// consonantVowelBoundaries returns a slice of 2-length strings that represent the unique boundaries between consonant
// groups and vowel groups in the word. For example, "apple" would return ["le"] because it is the only consonant
// group ("ppl") followed by a vowel group ("e"), "street" would return ["re"], "banana" would return ["ba", "na"], and
// "beautiful" would return ["be", "ti", "fu"].
func consonantVowelBoundaries(word string) []string {
	vowels := "aeiouy"
	var boundaries []string
	uniqueBoundaries := make(map[string]bool)

	for i := 0; i < len(word)-1; i++ {
		currentChar := word[i]
		nextChar := word[i+1]

		currentIsVowel := strings.ContainsRune(vowels, rune(currentChar))
		nextIsVowel := strings.ContainsRune(vowels, rune(nextChar))

		if !currentIsVowel && nextIsVowel {
			boundary := string(currentChar) + string(nextChar)
			if !uniqueBoundaries[boundary] {
				boundaries = append(boundaries, boundary)
				uniqueBoundaries[boundary] = true
			}
		}
	}

	return boundaries
}

func PrefixesAndSuffixes(words map[string]int64) (map[string]bool, map[string]bool) {
	prefixes := make(map[string]int64)
	suffixes := make(map[string]int64)

	for word, count := range words {
		length := len(word)
		if length >= 3 {
			prefix := word[:2]
			suffix := word[length-2:]
			prefixes[prefix] += count
			suffixes[suffix] += count
		}
		if length >= 4 {
			prefix := word[:3]
			suffix := word[length-3:]
			prefixes[prefix] += count
			suffixes[suffix] += count
		}
	}

	// Filter prefixes and suffixes based on the specified thresholds
	filteredPrefixes := make(map[string]bool)
	filteredSuffixes := make(map[string]bool)

	for prefix, _ := range prefixes {
		filteredPrefixes[prefix] = true
	}

	for suffix, _ := range suffixes {
		filteredSuffixes[suffix] = true
	}

	return filteredPrefixes, filteredSuffixes
}

func Middles(words map[string]int64) map[string]bool {
	middleSubstrings := make(map[string]int64)

	for word, count := range words {
		length := len(word)
		if length >= 4 {
			for i := 1; i < length-2; i++ {
				substr := word[i : i+2]
				middleSubstrings[substr] += count
			}
		}
		if length >= 5 {
			for i := 1; i < length-3; i++ {
				substr := word[i : i+3]
				middleSubstrings[substr] += count
			}
		}
		if length >= 6 {
			for i := 1; i < length-4; i++ {
				substr := word[i : i+4]
				middleSubstrings[substr] += count
			}
		}
	}

	filteredMiddles := make(map[string]bool)

	for middle, _ := range middleSubstrings {
		filteredMiddles[middle] = true
	}

	return filteredMiddles
}

func AntiPrefixes(words map[string]int64) map[string]bool {
	// Create a slice of all possible 2 and 3 letter combinations
	var combinations []string
	letters := "abcdefghijklmnopqrstuvwxyz"

	// Generate 2-letter combinations
	for _, first := range letters {
		for _, second := range letters {
			combinations = append(combinations, string([]rune{first, second}))
		}
	}

	// Generate 3-letter combinations
	for _, first := range letters {
		for _, second := range letters {
			for _, third := range letters {
				combinations = append(combinations, string([]rune{first, second, third}))
			}
		}
	}

	// Create a map to store the prefixes that occur in words
	prefixes := make(map[string]bool)

	// Iterate through the words and mark the prefixes that occur
	for word := range words {
		length := len(word)
		if length >= 3 {
			prefixes[word[:2]] = true
		}
		if length >= 4 {
			prefixes[word[:3]] = true
		}
	}

	// Eliminate any 2 or 3 letter prefixes that occur in words
	antiPrefixes := make(map[string]bool)
	for _, combination := range combinations {
		if !prefixes[combination] {
			antiPrefixes[combination] = true
		}
	}

	return antiPrefixes
}

func AntiSuffixes(words map[string]int64) map[string]bool {
	// Create a slice of all possible 2 and 3 letter combinations
	var combinations []string
	letters := "abcdefghijklmnopqrstuvwxyz"

	// Generate 2-letter combinations
	for _, first := range letters {
		for _, second := range letters {
			combinations = append(combinations, string([]rune{first, second}))
		}
	}

	// Generate 3-letter combinations
	for _, first := range letters {
		for _, second := range letters {
			for _, third := range letters {
				combinations = append(combinations, string([]rune{first, second, third}))
			}
		}
	}

	// Create a map to store the suffixes that occur in words
	suffixes := make(map[string]bool)

	// Iterate through the words and mark the suffixes that occur
	for word := range words {
		length := len(word)
		if length >= 3 {
			suffixes[word[length-2:]] = true
		}
		if length >= 4 {
			suffixes[word[length-3:]] = true
		}
	}

	// Eliminate any 2 or 3 letter suffixes that occur in words
	antiSuffixes := make(map[string]bool)
	for _, combination := range combinations {
		if !suffixes[combination] {
			antiSuffixes[combination] = true
		}
	}

	return antiSuffixes
}

func AntiMiddles(words map[string]int64) map[string]bool {
	// Create a slice of all possible 2 and 3 letter combinations
	var combinations []string
	letters := "abcdefghijklmnopqrstuvwxyz"

	// Generate 2-letter combinations
	for _, first := range letters {
		for _, second := range letters {
			combinations = append(combinations, string([]rune{first, second}))
		}
	}

	// Generate 3-letter combinations
	for _, first := range letters {
		for _, second := range letters {
			for _, third := range letters {
				combinations = append(combinations, string([]rune{first, second, third}))
			}
		}
	}

	// Create a map to store the middles that occur in words
	middles := make(map[string]bool)

	// Iterate through the words and mark the middles that occur
	for word := range words {
		length := len(word)
		if length >= 4 {
			for i := 1; i <= length-3; i++ {
				middles[word[i:i+2]] = true
			}
		}
		if length >= 5 {
			for i := 1; i <= length-4; i++ {
				middles[word[i:i+3]] = true
			}
		}
	}

	// Eliminate any 2 or 3 letter middles that occur in words
	antiMiddles := make(map[string]bool)
	for _, combination := range combinations {
		if !middles[combination] {
			antiMiddles[combination] = true
		}
	}

	return antiMiddles
}

func WriteCSV(partialWords map[string]int64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Partial Word", "Count"}); err != nil {
		return err
	}

	// Write data
	for k, v := range partialWords {
		if err := writer.Write([]string{k, strconv.FormatInt(v, 10)}); err != nil {
			return err
		}
	}

	return nil
}
