package sifo

import (
	"fmt"
	"math/rand"
	"time"
)

type Cipher map[string]string

const (
	firstThresholdFactor  = 6.5
	secondThresholdFactor = 2.0
	largeVariationsAfter  = 200

	quote = "Here's to the crazy ones. The misfits. The rebels. The troublemakers. The round pegs in the square holes."
)

var restarts int

func FindBestCipher(words map[string]int64, iterations int) Cipher {
	var bestCipher Cipher
	var maxHighScore float64

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	gts := giants(words)
	restarts++

	minGiantScore := gts[0].score
	for _, gt := range gts {
		fmt.Printf("%d. Giant %s: %.4f\n", restarts, gt.name, gt.score)
		fmt.Printf("  %s\n", encode(quote, gt.cipher))
		if gt.score < minGiantScore {
			minGiantScore = gt.score
		}
	}

	// Elastic concept: Avoids the problem of getting stuck on good prospects that may not have the highest potential.
	// Testing shows that many more promising ciphers can be varied to excellence and potentially greatness by not just
	// getting stuck on good candidates and comparing newly sprouted candidates against the more advanced candidates too
	// early. It takes the hardcoded approach but makes it more dynamic so it updates itself. Using this concept, I was
	// nearly able to replicate the findings of 2 million iterations in a fraction of the iterations. Essentially, it is
	// like having good performing managers keeping promising employees from rising, when they may be better, eventually,
	// than the current top performers.
	//
	// Threshold approach: Completely random to get past the first threshold, then vary it to find its potential.
	// Additional ideas: If it doesn't reach a second threshold in a certain number of iterations, then fall back to
	// a new completely random cipher at the first threshold.

	i := 0
	for {
		cipher := generateRandomCipherSimple(r) // call first to avoid bestCipher being empty
		highScore := CountHighScore(words, cipher, false)
		i++

		if highScore > maxHighScore {
			maxHighScore = highScore
			fmt.Printf("%d. New high score (pre-first threshold), %.4f: %d iterations\n", restarts, maxHighScore, i)
		}

		if maxHighScore > (minGiantScore / firstThresholdFactor) {
			fmt.Printf("%d. First threshold reached (%.4f > %.4f): %d iterations\n", restarts, maxHighScore, (minGiantScore / firstThresholdFactor), i)
			bestCipher = cipher
			maxHighScore = highScore
			break
		}
	}

	secondThresholdPassed := false
	thirdThresholdPassed := false
	iterationsSinceLastHighScore := 0

	for i := 0; i < iterations; i++ {
		tryCipher := bestCipher

		/*
			if thirdThresholdPassed || iterationsSinceLastHighScore > 4000 {
				switch r.Intn(6) {
				case 0:
					tryCipher = kindCatsCipher()
				case 1:
					tryCipher = ableEchoCipher()
				case 2:
					tryCipher = lazyFiveCipher()
				case 3:
					tryCipher = parisMilanCipher()
				case 4, 5:
					tryCipher = bestCipher
				}
			}
		*/

		variations := 1 + r.Intn(2)
		if iterationsSinceLastHighScore > largeVariationsAfter {
			variations = 1 + r.Intn(3)
		}

		cipher := varyCipher(tryCipher, r, variations)
		highScore := CountHighScore(words, cipher, false)
		if highScore > maxHighScore {
			isAGiant := false
			for _, gt := range gts {
				if equal(gt.cipher, cipher) {
					isAGiant = true
					break
				}
			}

			if !isAGiant {
				maxHighScore = highScore
				bestCipher = cipher
				iterationsSinceLastHighScore = 0

				thresh := "first"
				if secondThresholdPassed {
					thresh = "second"
				}
				if thirdThresholdPassed {
					thresh = "third"
				}
				fmt.Printf("%d. New high score (%s threshold), %.4f: %d iterations\n", restarts, thresh, maxHighScore, i)
			}
		}

		iterationsSinceLastHighScore++

		if maxHighScore > (minGiantScore/secondThresholdFactor) && !secondThresholdPassed {
			secondThresholdPassed = true
			fmt.Printf("%d. Second threshold reached (%.4f > %.4f): %d iterations\n", restarts, maxHighScore, minGiantScore/secondThresholdFactor, i)
		}

		if !secondThresholdPassed && i > (iterations/3) {
			fmt.Printf("%d. Second threshold not reached (%.4f < %.4f) after %d iterations, restarting...\n", restarts, maxHighScore, minGiantScore/secondThresholdFactor, i)
			return FindBestCipher(words, iterations)
		}

		if maxHighScore > minGiantScore && !thirdThresholdPassed {
			thirdThresholdPassed = true
			fmt.Printf("%d. Third threshold reached (%.4f > %.4f): %d iterations\n", restarts, maxHighScore, minGiantScore, i)
		}
	}

	beatAGiant := false
	isAGiant := false
	giantFound := ""

	for _, gt := range gts {
		if equal(gt.cipher, bestCipher) {
			isAGiant = true
			giantFound = gt.name
			break
		}
		if maxHighScore > gt.score {
			fmt.Printf("%d. Beat %s: %.4f > %.4f\n", restarts, gt.name, maxHighScore, gt.score)
			beatAGiant = true
		}
	}

	if isAGiant {
		fmt.Printf("%d. Is a giant (duplicate of %s), restarting...\n", restarts, giantFound)
		return FindBestCipher(words, iterations)
	}

	if !beatAGiant {
		fmt.Printf("%d. Didn't beat a giant after %d iterations, restarting...\n", restarts, iterations)
		return FindBestCipher(words, iterations)
	}

	fmt.Printf("%d. High score: %.4f\n", restarts, maxHighScore)
	return bestCipher
}

func CountHighScore(words map[string]int64, cipher Cipher, output bool) float64 {
	var count int64
	var score float64
	for word, ogScore := range words {
		encodedWord := encodeWord(word, cipher)
		if encScore, ok := words[encodedWord]; ok {
			if output {
				fmt.Printf("%s -> %s\n", word, encodedWord)
			}
			count++
			score = score + (float64(ogScore+encScore) / 2)
		}
	}
	if output {
		fmt.Printf("Count: %d\n", count)
		fmt.Printf("Score: %.4f\n", score)
	}
	return (float64(count) * score) / 25463612790
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

func giants(words map[string]int64) []Giant {
	giantsList := []Giant{
		{"KindCats", kindCatsCipher(), CountHighScore(words, kindCatsCipher(), false)},
		{"KindCats2", kindCats2Cipher(), CountHighScore(words, kindCats2Cipher(), false)},
		{"AbleEcho", ableEchoCipher(), CountHighScore(words, ableEchoCipher(), false)},
		{"LazyFive", lazyFiveCipher(), CountHighScore(words, lazyFiveCipher(), false)},
		{"ParisMilan", parisMilanCipher(), CountHighScore(words, parisMilanCipher(), false)},
		{"ParisMilan2", parisMilan2Cipher(), CountHighScore(words, parisMilan2Cipher(), false)},
		{"ParisMilan3", parisMilan3Cipher(), CountHighScore(words, parisMilan3Cipher(), false)},
		{"ParisMilan4", parisMilan4Cipher(), CountHighScore(words, parisMilan4Cipher(), false)},
		{"ParisMilan5", parisMilan5Cipher(), CountHighScore(words, parisMilan5Cipher(), false)},
		{"ParisMilan6", parisMilan6Cipher(), CountHighScore(words, parisMilan6Cipher(), false)},
		{"ParisMilan7", parisMilan7Cipher(), CountHighScore(words, parisMilan7Cipher(), false)},
	}

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

func kindCatsCipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "p"
	cipher["c"] = "j"
	cipher["d"] = "s"
	cipher["e"] = "o"
	cipher["f"] = "h"
	cipher["g"] = "y"
	cipher["h"] = "f"
	cipher["i"] = "a"
	cipher["j"] = "q"
	cipher["k"] = "c"
	cipher["l"] = "w"
	cipher["m"] = "r"
	cipher["n"] = "t"
	cipher["o"] = "e"
	cipher["p"] = "m"
	cipher["q"] = "z"
	cipher["r"] = "l"
	cipher["s"] = "n"
	cipher["t"] = "d"
	cipher["u"] = "k"
	cipher["v"] = "u"
	cipher["w"] = "b"
	cipher["x"] = "v"
	cipher["y"] = "g"
	cipher["z"] = "x"
	return cipher
}

func kindCats2Cipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "p"
	cipher["c"] = "j"
	cipher["d"] = "s"
	cipher["e"] = "o"
	cipher["f"] = "h"
	cipher["g"] = "y"
	cipher["h"] = "f"
	cipher["i"] = "a"
	cipher["j"] = "x"
	cipher["k"] = "c"
	cipher["l"] = "w"
	cipher["m"] = "r"
	cipher["n"] = "t"
	cipher["o"] = "e"
	cipher["p"] = "m"
	cipher["q"] = "z"
	cipher["r"] = "l"
	cipher["s"] = "n"
	cipher["t"] = "d"
	cipher["u"] = "k"
	cipher["v"] = "u"
	cipher["w"] = "b"
	cipher["x"] = "q"
	cipher["y"] = "g"
	cipher["z"] = "v"
	return cipher
}

func ableEchoCipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "e"
	cipher["b"] = "c"
	cipher["c"] = "g"
	cipher["d"] = "r"
	cipher["e"] = "o"
	cipher["f"] = "s"
	cipher["g"] = "b"
	cipher["h"] = "l"
	cipher["i"] = "a"
	cipher["j"] = "q"
	cipher["k"] = "v"
	cipher["l"] = "h"
	cipher["m"] = "y"
	cipher["n"] = "t"
	cipher["o"] = "i"
	cipher["p"] = "m"
	cipher["q"] = "z"
	cipher["r"] = "f"
	cipher["s"] = "d"
	cipher["t"] = "n"
	cipher["u"] = "k"
	cipher["v"] = "x"
	cipher["w"] = "p"
	cipher["x"] = "u"
	cipher["y"] = "w"
	cipher["z"] = "j"
	return cipher
}

func lazyFiveCipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "d"
	cipher["c"] = "w"
	cipher["d"] = "t"
	cipher["e"] = "o"
	cipher["f"] = "l"
	cipher["g"] = "m"
	cipher["h"] = "r"
	cipher["i"] = "u"
	cipher["j"] = "q"
	cipher["k"] = "j"
	cipher["l"] = "f"
	cipher["m"] = "n"
	cipher["n"] = "s"
	cipher["o"] = "a"
	cipher["p"] = "k"
	cipher["q"] = "x"
	cipher["r"] = "h"
	cipher["s"] = "b"
	cipher["t"] = "p"
	cipher["u"] = "c"
	cipher["v"] = "z"
	cipher["w"] = "g"
	cipher["x"] = "y"
	cipher["y"] = "e"
	cipher["z"] = "v"
	return cipher
}

func parisMilanCipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "p"
	cipher["c"] = "j"
	cipher["d"] = "s"
	cipher["e"] = "o"
	cipher["f"] = "h"
	cipher["g"] = "y"
	cipher["h"] = "f"
	cipher["i"] = "a"
	cipher["j"] = "x"
	cipher["k"] = "c"
	cipher["l"] = "w"
	cipher["m"] = "r"
	cipher["n"] = "t"
	cipher["o"] = "e"
	cipher["p"] = "m"
	cipher["q"] = "z"
	cipher["r"] = "l"
	cipher["s"] = "n"
	cipher["t"] = "d"
	cipher["u"] = "k"
	cipher["v"] = "u"
	cipher["w"] = "b"
	cipher["x"] = "v"
	cipher["y"] = "g"
	cipher["z"] = "q"
	return cipher
}

func parisMilan2Cipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "p"
	cipher["c"] = "j"
	cipher["d"] = "s"
	cipher["e"] = "o"
	cipher["f"] = "h"
	cipher["g"] = "y"
	cipher["h"] = "f"
	cipher["i"] = "a"
	cipher["j"] = "z"
	cipher["k"] = "c"
	cipher["l"] = "w"
	cipher["m"] = "r"
	cipher["n"] = "t"
	cipher["o"] = "e"
	cipher["p"] = "m"
	cipher["q"] = "x"
	cipher["r"] = "l"
	cipher["s"] = "n"
	cipher["t"] = "d"
	cipher["u"] = "k"
	cipher["v"] = "u"
	cipher["w"] = "b"
	cipher["x"] = "v"
	cipher["y"] = "g"
	cipher["z"] = "q"
	return cipher
}

func parisMilan3Cipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "p"
	cipher["c"] = "j"
	cipher["d"] = "s"
	cipher["e"] = "o"
	cipher["f"] = "h"
	cipher["g"] = "y"
	cipher["h"] = "f"
	cipher["i"] = "a"
	cipher["j"] = "x"
	cipher["k"] = "c"
	cipher["l"] = "w"
	cipher["m"] = "r"
	cipher["n"] = "t"
	cipher["o"] = "e"
	cipher["p"] = "m"
	cipher["q"] = "z"
	cipher["r"] = "l"
	cipher["s"] = "n"
	cipher["t"] = "d"
	cipher["u"] = "k"
	cipher["v"] = "q"
	cipher["w"] = "b"
	cipher["x"] = "v"
	cipher["y"] = "g"
	cipher["z"] = "u"
	return cipher
}

func parisMilan4Cipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "p"
	cipher["c"] = "j"
	cipher["d"] = "s"
	cipher["e"] = "o"
	cipher["f"] = "h"
	cipher["g"] = "y"
	cipher["h"] = "f"
	cipher["i"] = "a"
	cipher["j"] = "x"
	cipher["k"] = "c"
	cipher["l"] = "w"
	cipher["m"] = "r"
	cipher["n"] = "t"
	cipher["o"] = "e"
	cipher["p"] = "m"
	cipher["q"] = "u"
	cipher["r"] = "l"
	cipher["s"] = "n"
	cipher["t"] = "d"
	cipher["u"] = "k"
	cipher["v"] = "z"
	cipher["w"] = "b"
	cipher["x"] = "v"
	cipher["y"] = "g"
	cipher["z"] = "q"
	return cipher
}

func parisMilan5Cipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "p"
	cipher["c"] = "j"
	cipher["d"] = "s"
	cipher["e"] = "o"
	cipher["f"] = "h"
	cipher["g"] = "y"
	cipher["h"] = "f"
	cipher["i"] = "a"
	cipher["j"] = "q"
	cipher["k"] = "c"
	cipher["l"] = "w"
	cipher["m"] = "r"
	cipher["n"] = "t"
	cipher["o"] = "e"
	cipher["p"] = "m"
	cipher["q"] = "u"
	cipher["r"] = "l"
	cipher["s"] = "n"
	cipher["t"] = "d"
	cipher["u"] = "k"
	cipher["v"] = "z"
	cipher["w"] = "b"
	cipher["x"] = "v"
	cipher["y"] = "g"
	cipher["z"] = "x"
	return cipher
}

func parisMilan6Cipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "p"
	cipher["c"] = "j"
	cipher["d"] = "s"
	cipher["e"] = "o"
	cipher["f"] = "h"
	cipher["g"] = "y"
	cipher["h"] = "f"
	cipher["i"] = "a"
	cipher["j"] = "u"
	cipher["k"] = "c"
	cipher["l"] = "w"
	cipher["m"] = "r"
	cipher["n"] = "t"
	cipher["o"] = "e"
	cipher["p"] = "m"
	cipher["q"] = "z"
	cipher["r"] = "l"
	cipher["s"] = "n"
	cipher["t"] = "d"
	cipher["u"] = "k"
	cipher["v"] = "q"
	cipher["w"] = "b"
	cipher["x"] = "v"
	cipher["y"] = "g"
	cipher["z"] = "x"
	return cipher
}

func parisMilan7Cipher() Cipher {
	cipher := make(Cipher)
	cipher["a"] = "i"
	cipher["b"] = "p"
	cipher["c"] = "j"
	cipher["d"] = "s"
	cipher["e"] = "o"
	cipher["f"] = "h"
	cipher["g"] = "y"
	cipher["h"] = "f"
	cipher["i"] = "a"
	cipher["j"] = "q"
	cipher["k"] = "c"
	cipher["l"] = "w"
	cipher["m"] = "r"
	cipher["n"] = "t"
	cipher["o"] = "e"
	cipher["p"] = "m"
	cipher["q"] = "z"
	cipher["r"] = "l"
	cipher["s"] = "n"
	cipher["t"] = "d"
	cipher["u"] = "k"
	cipher["v"] = "x"
	cipher["w"] = "b"
	cipher["x"] = "v"
	cipher["y"] = "g"
	cipher["z"] = "u"
	return cipher
}
