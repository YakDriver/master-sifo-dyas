package sifo

import (
	"strings"
)

func encode(input string, cipher Cipher) string {
	var encoded strings.Builder
	for _, word := range strings.Fields(input) {
		encoded.WriteString(encodeWord(word, cipher))
		encoded.WriteRune(' ')
	}
	return encoded.String()
}

func encodeWord(word string, cipher Cipher) string {
	var encoded strings.Builder
	i := 0
	for i < len(word) {
		matched := false
		for length := len(word) - i; length > 0; length-- {
			substr := word[i : i+length]
			if encodedChars, ok := cipher[strings.ToLower(substr)]; ok {
				for j, char := range encodedChars {
					if i+j < len(word) && isUpper(word[i+j]) {
						encoded.WriteRune(toUpper(char))
					} else {
						encoded.WriteRune(char)
					}
				}
				i += length
				matched = true
				break
			}
		}
		if !matched {
			if isUpper(word[i]) {
				encoded.WriteRune(toUpper(rune(word[i])))
			} else {
				encoded.WriteRune(rune(word[i]))
			}
			i++
		}
	}
	return encoded.String()
}

func isUpper(r byte) bool {
	return r >= 'A' && r <= 'Z'
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 'a' + 'A'
	}
	return r
}

// decodeWord decodes a word using the given cipher from the end of the word to the beginning. In attempting to use
// multi-character ciphers, it will try to match the longest possible cipher first. This does not overcome the issue of
// ambiguous ciphers, but it does help to reduce the number of ambiguous ciphers.
func decodeWord(encodedWord string, cipher Cipher) string {
	// Create a reverse cipher map
	reverseCipher := make(map[string]string)
	for key, value := range cipher {
		reverseCipher[value] = key
	}

	var decoded string
	i := len(encodedWord)
	for i > 0 {
		matched := false
		for length := i; length > 0; length-- {
			substr := encodedWord[i-length : i]
			if decodedChars, ok := reverseCipher[substr]; ok {
				decoded = decodedChars + decoded
				i -= length
				matched = true
				break
			}
		}
		if !matched {
			decoded = string(encodedWord[i-1]) + decoded
			i--
		}
	}
	return decoded
}

func decodeWord2(encodedWord string, cipher Cipher) string {
	// Create a reverse cipher map
	reverseCipher := make(map[string]string)
	for key, value := range cipher {
		reverseCipher[value] = key
	}

	return encodeWord(encodedWord, reverseCipher)
}
