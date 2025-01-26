package main

import (
	"fmt"

	"github.com/YakDriver/master-sifo-dyas/sifo"
)

func main() {
	words := sifo.LoadWords("words.csv")
	prefixes, suffixes := sifo.PrefixesAndSuffixes(words)
	middles := sifo.Middles(words)

	fmt.Printf("Prefixes: %d\n", len(prefixes))
	fmt.Printf("Suffixes: %d\n", len(suffixes))
	fmt.Printf("Middles: %d\n", len(middles))

	antiPrefixes := sifo.AntiPrefixes(words)
	antiSuffixes := sifo.AntiSuffixes(words)
	antiMiddles := sifo.AntiMiddles(words)

	fmt.Printf("Anti Prefixes: %d\n", len(antiPrefixes))
	fmt.Printf("Anti Suffixes: %d\n", len(antiSuffixes))
	fmt.Printf("Anti Middles: %d\n", len(antiMiddles))

	wordPatterns := sifo.WordPatterns(words)
	vowelGroups := sifo.VowelGroups(words)
	consonantGroups := sifo.ConsonantGroups(words)

	fmt.Printf("Word Patterns: %d\n", len(wordPatterns))
	fmt.Printf("Vowel Groups: %d\n", len(vowelGroups))
	fmt.Printf("Consonant Groups: %d\n", len(consonantGroups))

	vowelConsonantBoundaries := sifo.VowelConsonantBoundaries(words)
	ConsonantVowelBoundaries := sifo.ConsonantVowelBoundaries(words)

	fmt.Printf("Vowel Consonant Boundaries: %d\n", len(vowelConsonantBoundaries))
	fmt.Printf("Consonant Vowel Boundaries: %d\n", len(ConsonantVowelBoundaries))

	//if 1 == 1 {
	//	return
	//}

	if 1 == 0 {
		partialWords := sifo.CreatePartialWordDictionary(words)
		if err := sifo.WriteCSV(partialWords, "partial_words.csv"); err != nil {
			fmt.Printf("Error writing to CSV: %v\n", err)
		} else {
			fmt.Println("Partial word dictionary written to partial_words.csv")
		}
	}

	dict := sifo.Dictionary{
		Words:                    words,
		Prefixes:                 prefixes,
		Suffixes:                 suffixes,
		Middles:                  middles,
		AntiPrefixes:             antiPrefixes,
		AntiSuffixes:             antiSuffixes,
		AntiMiddles:              antiMiddles,
		WordPatterns:             wordPatterns,
		VowelGroups:              vowelGroups,
		ConsonantGroups:          consonantGroups,
		VowelConsonantBoundaries: vowelConsonantBoundaries,
		ConsonantVowelBoundaries: ConsonantVowelBoundaries,
	}

	fmt.Printf("%s\n", sifo.Encode("Empower growth by nurturing strengths, guiding with patience, and leading with love—because what we cultivate in others, we cultivate in ourselves.", sifo.WarmHoldCipher()))

	//sifo.Score(dict, sifo.MoonPeerCipher(), true)

	bestCipher := sifo.FindBestCipher(dict, 10000)
	sifo.Score(dict, bestCipher, true)

	fmt.Println("Best Cipher:")
	for k, v := range bestCipher {
		fmt.Printf("%s -> %s\n", k, string(v))
	}
}
