package sifo

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

type Cipher map[string]string

const (
	firstThresholdFactor  = 2.7
	secondThresholdFactor = 1.2
	largeVariationsAfter  = 200

	quote = "Here's to the crazy ones. The misfits. The rebels. The troublemakers. The round pegs in the square holes."
)

type Strategy string

const (
	StrategyRandom  Strategy = "random"
	StrategyElastic Strategy = "elastic"
	StrategyGiant   Strategy = "giant"
)

type Threshold struct {
	score      float64
	iterations int
}

type Dictionary struct {
	Words                    map[string]int64
	Prefixes                 map[string]bool
	Suffixes                 map[string]bool
	Middles                  map[string]bool
	AntiPrefixes             map[string]bool
	AntiSuffixes             map[string]bool
	AntiMiddles              map[string]bool
	WordPatterns             map[string]bool // vowel/consonant patterns
	VowelGroups              map[string]bool
	ConsonantGroups          map[string]bool
	VowelConsonantBoundaries map[string]bool
	ConsonantVowelBoundaries map[string]bool
}

var restarts int

func FindBestCipher(dict Dictionary, iterations int) Cipher {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	gts := giants(dict)
	restarts++

	minGiantScore := gts[0].score
	for _, gt := range gts {
		fmt.Printf("%d. Giant %s: %.4f\n", restarts, gt.name, gt.score)
		fmt.Printf("  %s\n", Encode(quote, gt.cipher))
		if gt.score < minGiantScore {
			minGiantScore = gt.score
		}
	}

	var bestCipher Cipher

	bestCipher = generateRandomCipherSimple(r)

	objectiveAchieved := false

	for !objectiveAchieved {
		bestCipher, objectiveAchieved = iterationSearch(StrategyGiant, bestCipher, gts, dict, iterations, []Threshold{})
		if objectiveAchieved {
			break
		}
		bestCipher, objectiveAchieved = iterationSearch(StrategyRandom, generateRandomCipherSimple(r), gts, dict, iterations, []Threshold{{minGiantScore / firstThresholdFactor, 0}})
		if objectiveAchieved {
			bestCipher, objectiveAchieved = iterationSearch(StrategyElastic, bestCipher, gts, dict, iterations, []Threshold{
				{minGiantScore / secondThresholdFactor, iterations / 3},
				{minGiantScore, 0},
			})
		}
		restarts++
	}

	if objectiveAchieved {
		fmt.Printf("Objective achieved\n")
		return bestCipher
	}

	return bestCipher
}

// iterationSearch uses the strategy and returns the best cipher found and the high score.
// Whether iterations or thresholds are used depends on the strategy.
//
// Elastic means as long as thresholds are met, it will continue until the iterations are exhausted.
// In combination with random, it avoids the problem of getting stuck on good prospects that may not
// have the highest potential. Testing shows that many more promising ciphers can be varied to excellence
// and potentially greatness by not just getting stuck on good candidates and comparing newly sprouted
// candidates against the more advanced candidates too early. It takes the hardcoded approach but makes
// it more dynamic so it updates itself. Using this concept, I was nearly able to replicate the findings
// of 2 million iterations in a fraction of the iterations. Essentially, it is like having good performing
// managers keeping promising employees from rising, when they may be better, eventually, than the current
// top performers.
//
// Random approach means only the first threshold is used and goes until it is reached.
//
// Giants means using the giants as a reference point. If the high score is greater than all the giants,
// it will continue to vary the cipher. If the high score is less than or equal to a giant, it will vary
// the cipher based on the giant's cipher.
func iterationSearch(strat Strategy, bestCipher Cipher, gts []Giant, dict Dictionary, iterations int, thresholds []Threshold) (Cipher, bool) {
	var maxHighScore float64
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	thresholdsPassed := make([]bool, len(thresholds))
	curThreshold := 0

	maxHighScore = Score(dict, bestCipher, false)

	objectiveAchieved := false

	itsSinceHighScore := 0
outerLoop:
	for i := 0; i < iterations || strat == StrategyRandom || (iterations == 0 && curThreshold > len(thresholds)-1); i++ {
		tryCipher := bestCipher
		itsSinceHighScore++

		// by the end of each case, cipher should be ready to go
		switch strat {
		case StrategyRandom:
			if maxHighScore > thresholds[0].score {
				fmt.Printf("%d. First threshold reached (%.4f > %.4f): %d iterations\n", restarts, maxHighScore, thresholds[0].score, i)
				return bestCipher, true
			}
			tryCipher = generateRandomCipherSimple(r)
		case StrategyElastic:
			// returns in two ways: when the iterations are exhausted or if a threshold is not reached
			if curThreshold < len(thresholds) && maxHighScore > thresholds[curThreshold].score && !thresholdsPassed[curThreshold] {
				thresholdsPassed[curThreshold] = true
				fmt.Printf("%d. %s threshold reached (%.4f > %.4f): %d iterations\n", restarts, whichThreshold(curThreshold), maxHighScore, thresholds[curThreshold].score, i)

				curThreshold++

				if curThreshold > len(thresholds)-1 {
					objectiveAchieved = true
				}
			}

			// exit for not reaching threshold
			if curThreshold < (len(thresholds)-1) && maxHighScore < thresholds[curThreshold].score && i > thresholds[curThreshold].iterations {
				fmt.Printf("%d. %s threshold not reached (%.4f < %.4f) after %d iterations\n", restarts, whichThreshold(curThreshold), maxHighScore, thresholds[0].score, i)
				return bestCipher, false
			}

			tryCipher = varyCipher(bestCipher, r, 1+r.Intn(2))

			if itsSinceHighScore > 1800 {
				break outerLoop
			}
		case StrategyGiant:
			maxGiant := 0
			maxGiantScore := gts[0].score
			for i, gt := range gts {
				if maxGiantScore > gt.score {
					maxGiantScore = gt.score
					maxGiant = i
				}
			}

			if maxHighScore > maxGiantScore {
				objectiveAchieved = true
			}

			variations := 1 + r.Intn(2)
			if i%5 == 0 {
				variations = 1
			}

			if i%5 == 0 && maxHighScore > maxGiantScore {
				tryCipher = varyCipher(bestCipher, r, variations)
			} else {
				whichGiant := 0
				for {
					whichGiant = r.Intn(len(gts))
					if whichGiant != maxGiant {
						break
					}
				}
				tryCipher = varyCipher(gts[whichGiant].cipher, r, variations)
			}

			if itsSinceHighScore > 1000 {
				break outerLoop
			}
		}

		highScore := Score(dict, tryCipher, false)
		if int64(highScore) > int64(maxHighScore) {
			itsSinceHighScore = 0
			maxHighScore = highScore
			bestCipher = tryCipher

			id := ""
			switch strat {
			case StrategyRandom:
				id = fmt.Sprintf("(%s, %s threshold)", strat, whichThreshold(curThreshold))
			case StrategyElastic:
				id = fmt.Sprintf("(%s, %s threshold)", strat, whichThreshold(curThreshold))
			case StrategyGiant:
				id = fmt.Sprintf("(%s)", strat)
			}

			fmt.Printf("%s New high score, %.4f: %d iterations\n", id, maxHighScore, i)
		}

		if i%1000 == 0 && itsSinceHighScore > 500 {
			fmt.Printf("%d iterations, high score: %.4f, last high score was %d iterations ago\n", i, maxHighScore, itsSinceHighScore)
		}
	}

	if objectiveAchieved {
		isAGiant := false
		for _, gt := range gts {
			if equal(gt.cipher, bestCipher) {
				isAGiant = true
				break
			}
		}

		if isAGiant {
			objectiveAchieved = false
		}
	}

	return bestCipher, objectiveAchieved
}

func whichThreshold(i int) string {
	thresholds := []string{
		"first", "second", "third", "fourth", "fifth",
		"sixth", "seventh", "eighth", "ninth", "tenth",
	}

	if i >= 0 && i < len(thresholds) {
		return thresholds[i]
	}
	return ""
}

func Score(dict Dictionary, cipher Cipher, output bool) float64 {
	var score float64
	i := 0
	for word, ogOccurence := range dict.Words {
		i++
		encodedWord := encodeWord(word, cipher)
		if encOccurence, ok := dict.Words[encodedWord]; ok {
			s := float64(max(occurrenceScore(ogOccurence), occurrenceScore(encOccurence)))
			s *= 10
			score = score + s

			if output {
				fmt.Printf("%d. %s -> %s (score %.4f)\n", i, word, encodedWord, s)
			}
			continue
		}

		epc := englishPattern(encodedWord, dict)
		s := adjustEPC(epc) * float64(occurrenceScore(ogOccurence))
		score = score + s

		if output {
			//fmt.Printf("%d. %s -> %s (pattern match, score %.4f, epc %d)\n", i, word, encodedWord, s, epc)
		}
	}
	if output {
		fmt.Printf("Score: %.4f\n", score)
	}
	return score
}

// adjustEPC adjusts the English Pattern score. EPC can be between 0 and 8 and results in this mapping:
// 0	0.0
// 1	0.0
// 2	0.0
// 3	0.7
// 4	1.9
// 5	3.5
// 6	5.4
// 7	7.6
// 8	10.0
func adjustEPC(epc int) float64 {
	return math.Pow(max(float64(epc-4), 0), 1.5) * 1.3
}

// occurrenceScore takes the occurrences associated with a word and returns a score of 1, 2, or 3 based on the
// occurrence. 3 is ~95th percentile words, where s > 160000. 2 is ~80th where s > 40000. 1 is the rest.
func occurrenceScore(s int64) float64 {
	if s > 600000 {
		return 5
	} else if s > 160000 {
		return 2
	} else if s > 40000 {
		return 1.5
	} else {
		return 1
	}
}

// englishPattern checks the word for 3 things:
// 1. Its pattern (found with wordPattern()) matches a known pattern in dict.VowelConsonantPatterns
// 2. Each vowel group (found with vowelGroups()) in the word matches a known vowel group in dict.VowelGroups
// 3. Each consonant group (found with consonantGroups()) in the word matches a known consonant group in dict.ConsonantGroups
func englishPattern(word string, dict Dictionary) int {
	matches := 0
	pattern := wordPattern(word)
	if dict.WordPatterns[pattern] {
		matches++
	}

	vg := true
	vowelGroups := vowelGroups(word)
	for _, group := range vowelGroups {
		if !dict.VowelGroups[group] {
			vg = false
			break
		}
	}
	if vg {
		matches++
	}

	cg := true
	consonantGroups := consonantGroups(word)
	for _, group := range consonantGroups {
		if !dict.ConsonantGroups[group] {
			cg = false
			break
		}
	}
	if cg {
		matches++
	}

	vcb := true
	vcBoundaries := vowelConsonantBoundaries(word)
	for _, vc := range vcBoundaries {
		if !dict.VowelConsonantBoundaries[vc] {
			vcb = false
			break
		}
	}
	if vcb {
		matches++
	}

	cvb := true
	cvBoundaries := consonantVowelBoundaries(word)
	for _, cv := range cvBoundaries {
		if !dict.ConsonantVowelBoundaries[cv] {
			cvb = false
			break
		}
	}
	if cvb {
		matches++
	}

	if hasPrefix(word, dict.Prefixes) {
		matches++
	} else if hasPrefix(word, dict.AntiPrefixes) {
		matches--
	}

	if hasSuffix(word, dict.Suffixes) {
		matches++
	} else if hasPrefix(word, dict.AntiSuffixes) {
		matches--
	}

	if hasMostMiddles(word, dict.Middles) {
		matches++
	} else if hasMiddles(word, dict.AntiMiddles) {
		matches--
	}

	return matches
}

// patterns takes a word and returns a score based on the patterns it follows. A word gets points for each of the following patterns:
// 2 points - isCloseMatch, based on levenshtein distance (takes too long)
// 1 point - has a vowel
// 1 point - has one of the prefixes
// 1 point - has one of the suffixes
// 1 point - has one of the middles
func patterns(word string, dict Dictionary) int {
	score := 0

	// Check if the word is a close match to any word in the dictionary
	if 1 == 0 && isCloseMatch(word, dict.Words) {
		score += 2
	}

	// Check if the word has a vowel
	if hasVowel(word) {
		score += 1
	}

	// Check if the word has one of the prefixes
	if hasPrefix(word, dict.Prefixes) {
		score += 1
	} else if hasPrefix(word, dict.AntiPrefixes) {
		score -= 2
	}

	// Check if the word has one of the suffixes
	if hasSuffix(word, dict.Suffixes) {
		score += 1
	} else if hasSuffix(word, dict.AntiSuffixes) {
		score -= 2
	}

	// Check if the word has one of the middles
	if hasMiddles(word, dict.Middles) {
		score += 1
	}

	if hasMiddles(word, dict.AntiMiddles) {
		score -= 2
	}

	return score
}

func isCloseMatch(word string, words map[string]int64) bool {
	for w := range words {
		if levenshteinDistance(word, w) <= len(w)/3 { // 3 letter word, 1 off
			return true
		}
	}
	return false
}

func hasVowel(word string) bool {
	vowels := "aeiou"
	for _, char := range word {
		if strings.ContainsRune(vowels, char) {
			return true
		}
	}
	return false
}

// hasPrefix checks if a word has one of the prefixes. Prefixes are either 2 or 3 letters long. To have a 2-letter
// prefix, the word must be at least 3 letters long. To have a 3-letter prefix, the word must be at least 4 letters
// long.
func hasPrefix(word string, prefixes map[string]bool) bool {
	if len(word) >= 3 {
		prefix := word[:2]
		if _, exists := prefixes[prefix]; exists {
			return true
		}
	}
	if len(word) >= 4 {
		prefix := word[:3]
		if _, exists := prefixes[prefix]; exists {
			return true
		}
	}
	return false
}

// hasSuffix checks if a word has one of the suffixes. Suffixes are either 2 or 3 letters long. To have a 2-letter
// suffix, the word must be at least 3 letters long. To have a 3-letter suffix, the word must be at least 4 letters
// long.
func hasSuffix(word string, suffixes map[string]bool) bool {
	if len(word) >= 3 {
		suffix := word[len(word)-2:]
		if _, exists := suffixes[suffix]; exists {
			return true
		}
	}
	if len(word) >= 4 {
		suffix := word[len(word)-3:]
		if _, exists := suffixes[suffix]; exists {
			return true
		}
	}
	return false
}

// hasMiddles checks if all of a word's middles are keys in the middles map. Middles never include the first or last letter of a word. Middles
// are either 2, 3, or 4 letters long. To have a 2-letter middle, the word must be at least 4 letters long. To have a
// 3-letter middle, the word must be at least 5 letters long. To have a 4-letter middle, the word must be at least 6
// letters long.
func hasMiddles(word string, middles map[string]bool) bool {
	length := len(word)
	if length < 4 {
		return false
	}

	// Check 2-letter middles
	if length >= 4 {
		for i := 1; i <= length-3; i++ {
			middle := word[i : i+2]
			if _, exists := middles[middle]; !exists {
				return false
			}
		}
	}

	// Check 3-letter middles
	if length >= 5 {
		for i := 1; i <= length-4; i++ {
			middle := word[i : i+3]
			if _, exists := middles[middle]; !exists {
				return false
			}
		}
	}

	// Check 4-letter middles
	if length >= 6 {
		for i := 1; i <= length-5; i++ {
			middle := word[i : i+4]
			if _, exists := middles[middle]; !exists {
				return false
			}
		}
	}

	return true
}

// hasMostMiddles checks if half or more of a word's middles are keys in the middles map. Middles never include the
// first or last letter of a word. Middles are either 2, 3, or 4 letters long. To have a 2-letter middle, the word
// must be at least 4 letters long. To have a 3-letter middle, the word must be at least 5 letters long. To have a
// 4-letter middle, the word must be at least 6 letters long.
func hasMostMiddles(word string, middles map[string]bool) bool {
	length := len(word)
	if length < 4 {
		return false
	}

	totalMiddles := 0
	matchingMiddles := 0

	// Check 2-letter middles
	if length >= 4 {
		for i := 1; i <= length-3; i++ {
			middle := word[i : i+2]
			totalMiddles++
			if _, exists := middles[middle]; exists {
				matchingMiddles++
			}
		}
	}

	// Check 3-letter middles
	if length >= 5 {
		for i := 1; i <= length-4; i++ {
			middle := word[i : i+3]
			totalMiddles++
			if _, exists := middles[middle]; exists {
				matchingMiddles++
			}
		}
	}

	// Check 4-letter middles
	if length >= 6 {
		for i := 1; i <= length-5; i++ {
			middle := word[i : i+4]
			totalMiddles++
			if _, exists := middles[middle]; exists {
				matchingMiddles++
			}
		}
	}

	return float64(matchingMiddles) >= (float64(totalMiddles) * float64(0.60))
}

func levenshteinDistance(a, b string) int {
	// Implementation of the Levenshtein distance algorithm
	// This function calculates the number of single-character edits (insertions, deletions, or substitutions)
	// required to change one word into the other
	la, lb := len(a), len(b)
	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
	}
	for i := 0; i <= la; i++ {
		d[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		d[0][j] = j
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			d[i][j] = min(d[i-1][j]+1, min(d[i][j-1]+1, d[i-1][j-1]+cost))
		}
	}
	return d[la][lb]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func equal(a, b Cipher) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if b[k] != v {
			return false
		}
	}

	return true
}

func generateRandomCipherSimple(r *rand.Rand) Cipher {
	alphabet := []string{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
		"g",
		"h",
		"i",
		"j",
		"k",
		"l",
		"m",
		"n",
		"o",
		"p",
		"q",
		"r",
		"s",
		"t",
		"u",
		"v",
		"w",
		"x",
		"y",
		"z",
	}

	alphabetMap := shuffleMap(r, alphabet)

	cipher := make(Cipher)
	for k, v := range alphabetMap {
		cipher[k] = v
	}

	return cipher
}

// generateRandomCipher does not work because it generates invalid ciphers. No cipher key can be the prefix of another. This is
// known but I thought I could figure a way around it. I could not. For example, "a" and "an" cannot both be keys in the same
// cipher. This is because in decoding, you cannot know unambiguously whether two single letters are meant to be decoded as
// a single letter or as two letters. You can get around this by using a delimiter between characters but that is undesirable.
// This is kept for historical reasons.
func generateRandomCipher(r *rand.Rand) Cipher {
	twos := []string{
		"an",
		"ar",
		"as",
		"at",
		"ed",
		"en",
		"er",
		"ha",
		"he",
		"hi",
		"in",
		"is",
		"it",
		"nd",
		"of",
		"on",
		"or",
		"ou",
		"qu",
		"re",
		"th",
		"to",
	}
	alphabet := []string{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
		"g",
		"h",
		"i",
		"j",
		"k",
		"l",
		"m",
		"n",
		"o",
		"p",
		"q",
		"r",
		"s",
		"t",
		"u",
		"v",
		"w",
		"x",
		"y",
		"z",
	}

	twoMap := shuffleMap(r, twos)
	alphabetMap := shuffleMap(r, alphabet)

	cipher := make(Cipher)
	for k, v := range twoMap {
		cipher[k] = v
	}
	for k, v := range alphabetMap {
		cipher[k] = v
	}

	return cipher
}

// shuffleMap creates a new map with the same keys as values, but with the values shuffled. No key will
// reference itself as a value.
func shuffleMap(r *rand.Rand, ss []string) map[string]string {
	shuffled := make([]string, len(ss))
	copy(shuffled, ss)

	for {
		r.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

		// Check for self-references
		valid := true
		for i, s := range ss {
			if s == shuffled[i] {
				valid = false
				break
			}
		}

		if valid {
			break
		}
	}

	m := make(map[string]string)
	for i, s := range ss {
		m[s] = shuffled[i]
	}
	return m
}

func varyCipher(baseCipher Cipher, r *rand.Rand, variation int) Cipher {
	newCipher := make(Cipher)
	for k, v := range baseCipher {
		newCipher[k] = v
	}

	keys := make([]string, 0, len(newCipher))
	for k := range newCipher {
		keys = append(keys, k)
	}

	for i := 0; i < variation; i++ {
		// Randomly pick two keys to swap their values
		var key1, key2 string

		for {
			key1 = keys[r.Intn(len(keys))]
			key2 = keys[r.Intn(len(keys))]

			if key1 == newCipher[key2] || key2 == newCipher[key1] {
				continue
			}

			break
		}

		// Swap the values
		newCipher[key1], newCipher[key2] = newCipher[key2], newCipher[key1]
	}

	return newCipher
}

type Giant struct {
	name   string
	cipher Cipher
	score  float64
}

func giants(dict Dictionary) []Giant {
	fmt.Printf("Loading giants...\n")
	giantsList := []Giant{
		{"LonelyRemark", LonelyRemarkCipher(), Score(dict, LonelyRemarkCipher(), false)},
		{"MoonPeer", MoonPeerCipher(), Score(dict, MoonPeerCipher(), false)},
		{"WormHeld", WormHeldCipher(), Score(dict, WormHeldCipher(), false)},
		{"WarmHold", WarmHoldCipher(), Score(dict, WarmHoldCipher(), false)},
		{"WormHelp", WormHelpCipher(), Score(dict, WormHelpCipher(), false)},
	}

	fmt.Printf("Loaded %d giants\n", len(giantsList))

	fmt.Printf("Removing duplicate giants...\n")
	uniqueGiants := []Giant{}
	for i, giant1 := range giantsList {
		duplicate := false
		for j, giant2 := range giantsList {
			if i != j && equal(giant1.cipher, giant2.cipher) {
				fmt.Printf("Duplicate giant: %s == %s\n", giant1.name, giant2.name)
				duplicate = true
				break
			}
		}
		if !duplicate {
			uniqueGiants = append(uniqueGiants, giant1)
		}
	}

	return uniqueGiants
}

func LonelyRemarkCipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "o"
	cipher["b"] = "f"
	cipher["c"] = "d"
	cipher["d"] = "p"
	cipher["e"] = "a"
	cipher["f"] = "g"
	cipher["g"] = "b"
	cipher["h"] = "y"
	cipher["i"] = "u"
	cipher["j"] = "x"
	cipher["k"] = "v"
	cipher["l"] = "r"
	cipher["m"] = "n"
	cipher["n"] = "m"
	cipher["o"] = "e"
	cipher["p"] = "c"
	cipher["q"] = "j"
	cipher["r"] = "l"
	cipher["s"] = "t"
	cipher["t"] = "s"
	cipher["u"] = "i"
	cipher["v"] = "w"
	cipher["w"] = "h"
	cipher["x"] = "z"
	cipher["y"] = "k"
	cipher["z"] = "q"
	return cipher
}

func WarmHoldCipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "o"
	cipher["b"] = "w"
	cipher["c"] = "b"
	cipher["d"] = "n"
	cipher["e"] = "a"
	cipher["f"] = "g"
	cipher["g"] = "p"
	cipher["h"] = "y"
	cipher["i"] = "u"
	cipher["j"] = "x"
	cipher["k"] = "v"
	cipher["l"] = "r"
	cipher["m"] = "d"
	cipher["n"] = "m"
	cipher["o"] = "e"
	cipher["p"] = "f"
	cipher["q"] = "j"
	cipher["r"] = "l"
	cipher["s"] = "t"
	cipher["t"] = "s"
	cipher["u"] = "i"
	cipher["v"] = "c"
	cipher["w"] = "h"
	cipher["x"] = "z"
	cipher["y"] = "k"
	cipher["z"] = "q"
	return cipher
}

func WormHelpCipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "o"
	cipher["b"] = "w"
	cipher["c"] = "d"
	cipher["d"] = "m"
	cipher["e"] = "a"
	cipher["f"] = "g"
	cipher["g"] = "f"
	cipher["h"] = "y"
	cipher["i"] = "u"
	cipher["j"] = "x"
	cipher["k"] = "v"
	cipher["l"] = "n"
	cipher["m"] = "p"
	cipher["n"] = "r"
	cipher["o"] = "e"
	cipher["p"] = "b"
	cipher["q"] = "j"
	cipher["r"] = "l"
	cipher["s"] = "t"
	cipher["t"] = "s"
	cipher["u"] = "i"
	cipher["v"] = "c"
	cipher["w"] = "h"
	cipher["x"] = "z"
	cipher["y"] = "k"
	cipher["z"] = "q"
	return cipher
}

func MoonPeerCipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "o"
	cipher["b"] = "w"
	cipher["c"] = "d"
	cipher["d"] = "m"
	cipher["e"] = "a"
	cipher["f"] = "g"
	cipher["g"] = "f"
	cipher["h"] = "y"
	cipher["i"] = "u"
	cipher["j"] = "v"
	cipher["k"] = "j"
	cipher["l"] = "n"
	cipher["m"] = "p"
	cipher["n"] = "r"
	cipher["o"] = "e"
	cipher["p"] = "b"
	cipher["q"] = "x"
	cipher["r"] = "l"
	cipher["s"] = "t"
	cipher["t"] = "s"
	cipher["u"] = "i"
	cipher["v"] = "c"
	cipher["w"] = "h"
	cipher["x"] = "z"
	cipher["y"] = "k"
	cipher["z"] = "q"
	return cipher
}

func WormHeldCipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "o"
	cipher["b"] = "c"
	cipher["c"] = "b"
	cipher["d"] = "n"
	cipher["e"] = "a"
	cipher["f"] = "g"
	cipher["g"] = "p"
	cipher["h"] = "y"
	cipher["i"] = "u"
	cipher["j"] = "x"
	cipher["k"] = "v"
	cipher["l"] = "r"
	cipher["m"] = "d"
	cipher["n"] = "m"
	cipher["o"] = "e"
	cipher["p"] = "f"
	cipher["q"] = "j"
	cipher["r"] = "l"
	cipher["s"] = "t"
	cipher["t"] = "s"
	cipher["u"] = "i"
	cipher["v"] = "z"
	cipher["w"] = "h"
	cipher["x"] = "w"
	cipher["y"] = "k"
	cipher["z"] = "q"
	return cipher
}
